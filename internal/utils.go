package internal

import (
	"github.com/Masterminds/semver/v3"
)

func GetSemver(tag string) (*semver.Version, error) {
	return semver.NewVersion(tag)
}
