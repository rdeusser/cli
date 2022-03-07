package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rdeusser/cli/version"
)

func TestGetHumanVersion(t *testing.T) {
	version.Version = "0.1.0"

	t.Run("should contain the version prelease when set explicitly", func(t *testing.T) {
		version.GitDescribe = "v0.1.0"
		version.VersionPrerelease = "dev"

		assert.Equal(t, "v0.1.0-dev", version.GetHumanVersion())
	})

	t.Run("should set the version prerelease to 'dev' when not set", func(t *testing.T) {
		version.GitDescribe = ""
		version.VersionPrerelease = ""

		assert.Equal(t, "v0.1.0-dev", version.GetHumanVersion())
	})

	t.Run("should add the git commit when set", func(t *testing.T) {
		version.GitCommit = "acd3b9e"

		assert.Equal(t, "v0.1.0-dev (acd3b9e)", version.GetHumanVersion())
	})
}
