package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThatKifPlatformGeneratedSandboxInTmp(t *testing.T) {
	assert.Contains(t, NewSkrtPlatform().Sandbox, "/tmp/kif_")
}

func TestThatNilErrorDoesNotExitApplication(t *testing.T) {
	ExitOnError(nil)
}
