package cmd

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestThatKifPlatformGeneratedSandboxInTmp(t *testing.T) {
	kif, err := NewKifPlatform()
	assert.NoError(t, err)
	assert.Contains(t, kif.Sandbox, "/tmp/kif_")
}

func TestThatKifPlatformCreatedSandbox(t *testing.T) {
	kif, err := NewKifPlatform()
	assert.NoError(t, err)
	_, err = os.Stat(kif.Sandbox)
	assert.NoError(t, err, "/tmp/kif_")
}

func TestThatNilErrorDoesNotExitApplication(t *testing.T) {
	ExitOnError(nil)
}
