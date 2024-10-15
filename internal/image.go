package internal

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
)

func IsOfficialImage(image string) bool {
	return strings.Count(image, "/") == 0
}

func GetImageURL(image string) string {
	baseUrl := "https://registry.hub.docker.com/v2/repositories/"
	if IsOfficialImage(image) {
		baseUrl += "library/"
	}
	baseUrl += image + "/tags?page_size=100"
	return baseUrl
}

func FetchImageTags(image string) ([]string, error) {
	url := GetImageURL(image)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	type Tag struct {
		Name string `json:"name"`
	}

	var tagsResponse struct {
		Count   int    `json:"count"`
		Results []Tag  `json:"results"`
		Next    string `json:"next"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tagsResponse); err != nil {
		return nil, err
	}

	var tags []string
	for _, tag := range tagsResponse.Results {
		tags = append(tags, tag.Name)
	}

	return tags, nil
}

func semverInstance(tag string) (*semver.Version, error) {
	// Attempt to parse the tag as a semantic version
	v, err := semver.NewVersion(tag)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func FindLatestVersion(current *semver.Version, tags []string, major, minor, patch bool) (string, error) {
	type VersionTag struct {
		Version *semver.Version
		Tag     string
	}

	var versionTags []VersionTag

	// Collect valid semantic versions
	for _, tag := range tags {
		v, err := semver.NewVersion(tag)
		if err != nil {
			continue // Skip tags that are not valid semantic versions
		}

		// Suffix checks
		currentSplit := strings.Split(current.Original(), "-")
		tagSplit := strings.Split(tag, "-")
		if len(currentSplit) != len(tagSplit) {
			continue
		}

		if len(currentSplit) > 1 && len(tagSplit) > 1 {
			currentSuffix := strings.Split(current.Original(), "-")[1]
			tagSuffix := strings.Split(tag, "-")[1]
			if currentSuffix != tagSuffix {
				continue
			}
		}

		versionTags = append(versionTags, VersionTag{Version: v, Tag: tag})
	}

	if len(versionTags) == 0 {
		return "", nil
	}

	// Sort versions in ascending order
	sort.Slice(versionTags, func(i, j int) bool {
		return versionTags[i].Version.LessThan(versionTags[j].Version)
	})

	// Find the latest version according to the flags
	for i := len(versionTags) - 1; i >= 0; i-- {
		v := versionTags[i].Version
		tag := versionTags[i].Tag

		if v.Compare(current) <= 0 {
			continue // Skip versions not newer than current
		}

		accept := false
		if major && v.Major() != current.Major() {
			accept = true
		} else if minor && v.Major() == current.Major() && v.Minor() != current.Minor() {
			accept = true
		} else if patch && v.Major() == current.Major() && v.Minor() == current.Minor() && v.Patch() != current.Patch() {
			accept = true
		}

		if accept {
			return tag, nil
		}
	}

	return "", nil
}
