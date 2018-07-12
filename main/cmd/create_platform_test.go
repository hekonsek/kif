package cmd

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
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
