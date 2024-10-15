package internal

import (
	"testing"

	"github.com/Masterminds/semver/v3"
)

func TestFindLatestVersion(t *testing.T) {
	tests := []struct {
		name    string
		current string
		tags    []string
		major   bool
		minor   bool
		patch   bool
		want    string
		wantErr bool
	}{
		{
			name:    "latest patch version",
			current: "1.0.0",
			tags:    []string{"1.0.1", "1.0.2", "1.1.0"},
			major:   false,
			minor:   false,
			patch:   true,
			want:    "1.0.2",
			wantErr: false,
		},
		{
			name:    "latest minor version",
			current: "1.0.0",
			tags:    []string{"1.0.1", "1.1.0", "1.2.0"},
			major:   false,
			minor:   true,
			patch:   false,
			want:    "1.2.0",
			wantErr: false,
		},
		{
			name:    "latest major version",
			current: "1.0.0",
			tags:    []string{"1.0.1", "2.0.0", "3.0.0"},
			major:   true,
			minor:   false,
			patch:   false,
			want:    "3.0.0",
			wantErr: false,
		},
		{
			name:    "no newer version",
			current: "1.0.0",
			tags:    []string{"0.9.9", "1.0.0"},
			major:   true,
			minor:   true,
			patch:   true,
			want:    "",
			wantErr: false,
		},
		{
			name:    "invalid semver tags",
			current: "1.0.0",
			tags:    []string{"1.0.1", "invalid", "1.1.0"},
			major:   false,
			minor:   true,
			patch:   false,
			want:    "1.1.0",
			wantErr: false,
		},
		{
			name:    "suffix match",
			current: "1.0.0-beta",
			tags:    []string{"1.0.1-beta", "1.1.0-beta", "1.0.1"},
			major:   false,
			minor:   false,
			patch:   true,
			want:    "1.0.1-beta",
			wantErr: false,
		},
		{
			name:    "suffix mismatch",
			current: "1.0.0-beta",
			tags:    []string{"1.0.1-alpha", "1.1.0-beta", "1.0.1"},
			major:   false,
			minor:   false,
			patch:   true,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			current, err := semver.NewVersion(tt.current)
			if err != nil {
				t.Fatalf("invalid current version: %v", err)
			}
			result, err := FindLatestVersion(current, tt.tags, tt.major, tt.minor, tt.patch)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindLatestVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.want {
				t.Errorf("FindLatestVersion() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestIsOfficalImage(t *testing.T) {
	tests := []struct {
		name  string
		image string
		want  bool
	}{
		{
			name:  "official image",
			image: "nginx",
			want:  true,
		},
		{
			name:  "non-official image",
			image: "pytorch/pytorch",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsOfficialImage(tt.image)
			if result != tt.want {
				t.Errorf("IsOfficialImage() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestGetImageURL(t *testing.T) {
	tests := []struct {
		name  string
		image string
		want  string
	}{
		{
			name:  "official image",
			image: "nginx",
			want:  "https://registry.hub.docker.com/v2/repositories/library/nginx/tags?page_size=100",
		},
		{
			name:  "non-official image",
			image: "pytorch/pytorch",
			want:  "https://registry.hub.docker.com/v2/repositories/pytorch/pytorch/tags?page_size=100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetImageURL(tt.image)
			if result != tt.want {
				t.Errorf("GetImageURL() = %v, want %v", result, tt.want)
			}
		})
	}
}
