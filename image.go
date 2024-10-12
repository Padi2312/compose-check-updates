package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
)

func IsOfficalImage(image string) bool {
	return strings.Count(image, "/") == 0
}

func GetImageURL(image string) string {
	baseUrl := "https://registry.hub.docker.com/v2/repositories/"
	// Check if the image is official or non-official
	if IsOfficalImage(image) {
		// Offical image do not have a namespace
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

func GetSemver(tag string) (*semver.Version, error) {
	return semver.NewVersion(tag)
}

func FindLatestVersion(current *semver.Version, tags []string) (string, error) {
	var versions []*semver.Version
	var suffix string

	// Determine the suffix of the current version
	if idx := strings.Index(current.Original(), "-"); idx != -1 {
		suffix = current.Original()[idx:]
	}

	for _, tag := range tags {
		v, err := semver.NewVersion(tag)
		if err == nil { // Skip non-semver tags
			// Check if the tag has the same suffix as the current version
			if suffix == "" || strings.HasSuffix(tag, suffix) {
				versions = append(versions, v)
			}
		}
	}

	if len(versions) == 0 {
		return "", nil
	}

	sort.Sort(semver.Collection(versions))
	for i := len(versions) - 1; i >= 0; i-- {
		if current.Compare(versions[i]) == -1 {
			return versions[i].Original(), nil
		}
	}

	return "", nil
}
