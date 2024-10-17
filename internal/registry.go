package internal

import (
	"encoding/json"
	"net/http"
	"strings"
)

type IRegistry interface {
	FetchImageTags(image string) ([]string, error)
}

type Registry struct {
	url string
}

func NewRegistry(url string) *Registry {
	if url == "" {
		url = "https://registry.hub.docker.com/v2/repositories/"
	}
	return &Registry{url: url}
}

func (r *Registry) FetchImageTags(image string) ([]string, error) {
	var tags []string
	url := r.getImageURL(image)

	for {
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

		for _, tag := range tagsResponse.Results {
			tags = append(tags, tag.Name)
		}

		if tagsResponse.Next == "" {
			break
		}

		url = tagsResponse.Next
	}

	return tags, nil
}

func (r *Registry) isOfficialImage(image string) bool {
	return strings.Count(image, "/") == 0
}

func (r *Registry) getImageURL(image string) string {
	// Remove tags from image name
	image = strings.Split(image, ":")[0]
	baseUrl := r.url
	if !strings.HasSuffix(baseUrl, "/") {
		baseUrl += "/"
	}
	if r.isOfficialImage(image) {
		baseUrl += "library/"
	}
	baseUrl += image + "/tags?page_size=100"
	return baseUrl
}
