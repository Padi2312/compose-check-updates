package internal

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
)

type TestFindLatestVersionStruct struct {
	Current  string
	Tags     []string
	Major    bool
	Minor    bool
	Patch    bool
	Expected string
}

func TestFindLatestVersion(t *testing.T) {
	tests := []struct {
		name     string
		testData struct {
			Current  string
			Tags     []string
			Major    bool
			Minor    bool
			Patch    bool
			Expected string
		}
	}{
		{
			name: "patch update available",
			testData: TestFindLatestVersionStruct{
				Current:  "1.0.0",
				Tags:     []string{"1.0.1", "1.0.2", "1.1.0"},
				Major:    false,
				Minor:    false,
				Patch:    true,
				Expected: "1.0.2",
			},
		},
		{
			name: "minor update available",
			testData: TestFindLatestVersionStruct{
				Current:  "1.0.0",
				Tags:     []string{"1.0.1", "1.1.0", "1.2.0"},
				Major:    false,
				Minor:    true,
				Patch:    false,
				Expected: "1.2.0",
			},
		},
		{
			name: "major update available",
			testData: TestFindLatestVersionStruct{
				Current:  "1.0.0",
				Tags:     []string{"1.0.1", "1.1.0", "2.0.0", "3.0.0"},
				Major:    true,
				Minor:    false,
				Patch:    false,
				Expected: "3.0.0",
			},
		},
		{
			name: "major update available with minor and patch",
			testData: TestFindLatestVersionStruct{
				Current:  "1.0.0",
				Tags:     []string{"1.0.1", "1.1.0", "2.0.0", "3.0.0", "3.1.2"},
				Major:    true,
				Minor:    false,
				Patch:    false,
				Expected: "3.1.2",
			},
		},
		{
			name: "no update available",
			testData: TestFindLatestVersionStruct{
				Current:  "1.0.0",
				Tags:     []string{"0.9.9", "1.0.0"},
				Major:    true,
				Minor:    true,
				Patch:    true,
				Expected: "",
			},
		},
		{
			name: "prerelease patch update available",
			testData: TestFindLatestVersionStruct{
				Current:  "1.0.0-beta",
				Tags:     []string{"1.0.1-beta", "1.1.0-beta", "1.2.0"},
				Major:    false,
				Minor:    false,
				Patch:    true,
				Expected: "1.0.1-beta",
			},
		},
		{
			name: "prerelease patch update not available",
			testData: TestFindLatestVersionStruct{
				Current:  "1.0.0-beta",
				Tags:     []string{"1.0.1-alpha", "1.1.0-beta", "1.1.0"},
				Major:    false,
				Minor:    false,
				Patch:    true,
				Expected: "",
			},
		},
		{
			name: "major update with prerelease",
			testData: TestFindLatestVersionStruct{
				Current:  "1.0.0",
				Tags:     []string{"2.0.0-alpha", "2.0.0-beta", "2.0.0"},
				Major:    true,
				Minor:    false,
				Patch:    false,
				Expected: "2.0.0",
			},
		},
		{
			name: "minor update with prerelease",
			testData: TestFindLatestVersionStruct{
				Current:  "1.0.0",
				Tags:     []string{"1.1.0-alpha", "1.1.0-beta", "1.1.0"},
				Major:    false,
				Minor:    true,
				Patch:    false,
				Expected: "1.1.0",
			},
		},
		{
			name: "patch update with prerelease",
			testData: TestFindLatestVersionStruct{
				Current:  "1.0.0",
				Tags:     []string{"1.0.1-alpha", "1.0.1-beta", "1.0.1"},
				Major:    false,
				Minor:    false,
				Patch:    true,
				Expected: "1.0.1",
			},
		},
		{
			name: "huge version jump",
			testData: TestFindLatestVersionStruct{
				Current:  "1.0.0",
				Tags:     []string{"1.0.1", "1.1.0", "2.0.0", "3.0.0", "4.0.0", "5.0.0", "6.0.0", "7.0.0", "8.0.0", "9.0.0", "10.0.0"},
				Major:    true,
				Minor:    false,
				Patch:    false,
				Expected: "10.0.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			current, err := semver.NewVersion(tt.testData.Current)
			assert.NoError(t, err, "invalid current version")

			result := FindLatestVersion(current, tt.testData.Tags, tt.testData.Major, tt.testData.Minor, tt.testData.Patch)
			assert.Equal(t, tt.testData.Expected, result)
		})
	}
}

func TestSuffixMismatch(t *testing.T) {
	test := TestFindLatestVersionStruct{
		Current:  "1.0.0-beta",
		Tags:     []string{"1.0.1-alpha", "1.1.0-beta", "1.1.0"},
		Major:    false,
		Minor:    false,
		Patch:    true,
		Expected: "",
	}
	current, err := semver.NewVersion(test.Current)
	assert.NoError(t, err, "invalid current version")

	result := FindLatestVersion(current, test.Tags, test.Major, test.Minor, test.Patch)
	assert.Equal(t, test.Expected, result)
}
