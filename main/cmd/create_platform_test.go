package cmd

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"gopkg.in/yaml.v2"
)

func TestThatNilErrorDoesNotExitApplication(t *testing.T) {
	ExitOnError(nil)
}

func TestThatKifPlatformGeneratedSandboxInTmp(t *testing.T) {
	kif, err := NewKifPlatform()
	assert.NoError(t, err)
	assert.Contains(t, kif.Sandbox, "/tmp/kif_")
}

func TestThatDefaultKifPlatformConfigurationIsEmpty(t *testing.T) {
	kif, err := NewKifPlatform()
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{}, kif.Configuration)
}

func TestThatKifPlatformCreatedSandbox(t *testing.T) {
	kif, err := NewKifPlatform()
	assert.NoError(t, err)
	_, err = os.Stat(kif.Sandbox)
	assert.NoError(t, err)
}

func TestThatKifPlatformRenderedIssuer(t *testing.T) {
	kif, err := NewKifPlatform()
	assert.NoError(t, err)
	err = kif.RenderTemplate("templates/issuer-letsencrypt.yml")
	assert.NoError(t, err)
	_, err = os.Stat(kif.Sandbox + "/templates/issuer-letsencrypt.yml")
	assert.NoError(t, err)
}

func TestThatKifPlatformRenderedChartWithNameAndVersion(t *testing.T) {
	kif, err := NewKifPlatform()
	assert.NoError(t, err)
	kif.Configuration["Chart"] = map[string]interface{}{
		"Name":    "SomeName",
		"Version": "SomeVersion",
	}
	err = kif.RenderTemplate("Chart.yaml")
	assert.NoError(t, err)
	chart, err := ioutil.ReadFile(kif.Sandbox + "/Chart.yaml")
	assert.NoError(t, err)
	chartText := string(chart)
	assert.Contains(t, chartText, `name: "SomeName"`)
	assert.Contains(t, chartText, `version: "SomeVersion"`)
}

func TestThatKifPlatformRenderedExtraRequirements(t *testing.T) {
	kif, err := NewKifPlatform()
	assert.NoError(t, err)
	err = ioutil.WriteFile("/tmp/extra-requirements.yml", []byte(`dependencies:
- name: foo
  repository: https://kubernetes-charts.storage.googleapis.com/
  version: 0.0.0`), 0644)
	assert.NoError(t, err)
	err = kif.RenderRequirements("/tmp/extra-requirements.yml")
	assert.NoError(t, err)
	chart, err := ioutil.ReadFile(kif.Sandbox + "/requirements.yaml")
	assert.NoError(t, err)
	chartText := string(chart)
	assert.Contains(t, chartText, `name: foo`)
}

func TestThatKifPlatformRenderedExtraValues(t *testing.T) {
	kif, err := NewKifPlatform()
	assert.NoError(t, err)
	err = ioutil.WriteFile("/tmp/extra-values.yml", []byte(`prometheus:
  foo: bar`), 0644)
	assert.NoError(t, err)
	err = kif.RenderValues("/tmp/extra-values.yml")
	assert.NoError(t, err)
	chart, err := ioutil.ReadFile(kif.Sandbox + "/values.yml")
	assert.NoError(t, err)
	generatedValues := map[string]map[string]interface{}{}
	err = yaml.Unmarshal(chart, &generatedValues)
	assert.NoError(t, err)
	assert.Equal(t, generatedValues["prometheus"]["foo"], "bar")
}

func TestThatKifPlatformExtraValuesMergePreservedExistingValues(t *testing.T) {
	kif, err := NewKifPlatform()
	assert.NoError(t, err)
	err = ioutil.WriteFile("/tmp/extra-values.yml", []byte(`prometheus:
  foo: bar`), 0644)
	assert.NoError(t, err)
	err = kif.RenderValues("/tmp/extra-values.yml")
	assert.NoError(t, err)
	chart, err := ioutil.ReadFile(kif.Sandbox + "/values.yml")
	assert.NoError(t, err)
	generatedValues := map[string]map[string]interface{}{}
	err = yaml.Unmarshal(chart, &generatedValues)
	assert.NoError(t, err)
	assert.NotNil(t, generatedValues["prometheus"]["alertmanager"])
}