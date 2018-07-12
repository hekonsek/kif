package cmd

import (
	"github.com/stretchr/testify/assert"
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

func TestThatKifPlatformCreatedSandbox(t *testing.T) {
	kif, err := NewKifPlatform()
	assert.NoError(t, err)
	_, err = os.Stat(kif.Sandbox)
	assert.NoError(t, err)
}

func TestThatKifPlatformRenderedIssuer(t *testing.T) {
	kif, err := NewKifPlatform()
	assert.NoError(t, err)
	config := map[string]interface{}{}
	err = kif.RenderTemplate("templates/issuer-letsencrypt", config)
	assert.NoError(t, err)
	_, err = os.Stat(kif.Sandbox + "/templates/issuer-letsencrypt.yml")
	assert.NoError(t, err)
}
