package internal

import (
	"regexp"
	"sort"

	"github.com/Masterminds/semver/v3"
)

func FindLatestVersion(current *semver.Version, tags []string, major, minor, patch bool) string {
	if major {
		minor = true
		patch = true
	}
	if minor {
		patch = true
	}

	type VersionTag struct {
		Version *semver.Version
		Tag     string
	}
	var versionTags []VersionTag

	// Collect valid semantic versions
	for _, tag := range tags {
		// Filter out invalid semantic versions
		if !isValidSemver(tag) {
			continue
		}

		// Attempt to parse the tag as a semantic version to compare it later easily
		v, err := semver.NewVersion(tag)
		if err != nil {
			continue
		}
		versionTags = append(versionTags, VersionTag{Version: v, Tag: tag})
	}

	if len(versionTags) == 0 {
		return ""
	}

	// Sort versions in descending order
	// This is necessary to find the latest version
	sort.Slice(versionTags, func(i, j int) bool {
		return versionTags[i].Version.GreaterThan(versionTags[j].Version)
	})

	for _, vt := range versionTags {
		v := vt.Version
		tag := vt.Tag

		// Skip versions not newer than current
		if v.LessThanEqual(current) || v.Prerelease() != current.Prerelease() {
			continue
		}

		accept := false
		if major && v.Major() > current.Major() {
			accept = true
		} else if minor && isEqualMajor(v, current) && v.Minor() > current.Minor() {
			accept = true
		} else if patch && isEqualMajor(v, current) && isEqualMinor(v, current) && v.Patch() > current.Patch() {
			accept = true
		}

		if accept {
			return tag
		}
	}

	return ""
}

func isValidSemver(tag string) bool {
	regex := regexp.MustCompile(`^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	matches := regex.FindStringSubmatch(tag)
	return len(matches) > 0
}

func isEqualMajor(current, tag *semver.Version) bool {
	return current.Major() == tag.Major()
}

func isEqualMinor(current, tag *semver.Version) bool {
	return current.Minor() == tag.Minor()
}
