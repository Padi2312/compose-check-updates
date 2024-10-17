package internal

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUpdateInfos(t *testing.T) {
	tests := []struct {
		name     string
		fileData string
		expected []UpdateInfo
	}{
		{
			name: "Single image",
			fileData: `
image: library/ubuntu:18.04.0
`,
			expected: []UpdateInfo{
				{
					RawLine:       "image: library/ubuntu:18.04.0",
					FullImageName: "library/ubuntu:18.04.0",
					ImageName:     "library/ubuntu",
					CurrentTag:    "18.04.0",
				},
			},
		},
		{
			name: "Multiple images",
			fileData: `
image: library/ubuntu:18.04.0
image: library/nginx:1.19.0
`,
			expected: []UpdateInfo{
				{
					RawLine:       "image: library/ubuntu:18.04.0",
					FullImageName: "library/ubuntu:18.04.0",
					ImageName:     "library/ubuntu",
					CurrentTag:    "18.04.0",
				},
				{
					RawLine:       "image: library/nginx:1.19.0",
					FullImageName: "library/nginx:1.19.0",
					ImageName:     "library/nginx",
					CurrentTag:    "1.19.0",
				},
			},
		},
		{
			name: "Duplicate images",
			fileData: `
image: library/ubuntu:18.04.0
image: library/ubuntu:18.04.0
`,
			expected: []UpdateInfo{
				{
					RawLine:       "image: library/ubuntu:18.04.0",
					FullImageName: "library/ubuntu:18.04.0",
					ImageName:     "library/ubuntu",
					CurrentTag:    "18.04.0",
				},
			},
		},
		{
			name: "No tag",
			fileData: `
image: library/ubuntu
`,
			expected: []UpdateInfo{
				{
					RawLine:       "image: library/ubuntu",
					FullImageName: "library/ubuntu",
					ImageName:     "library/ubuntu",
					CurrentTag:    "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file with the test data
			file, err := os.CreateTemp("", "testfile.yaml")
			assert.NoError(t, err)
			defer os.Remove(file.Name())

			_, err = file.WriteString(tt.fileData)
			assert.NoError(t, err)
			file.Close()

			// Update the expected FilePath to match the temporary file name
			for i := range tt.expected {
				tt.expected[i].FilePath = file.Name()
			}

			// Create an UpdateChecker instance
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"count": 2, "results": [
					{"name": "1.18.0"},
					{"name": "1.18.1"},
					{"name": "1.19.0"},
					{"name": "1.20.0"}
				],"next": null}`))
			}))
			defer server.Close()

			registry := NewRegistry(server.URL)
			updateChecker := NewUpdateChecker(file.Name(), registry)

			// Call createUpdateInfos
			updateInfos, err := updateChecker.createUpdateInfos()
			assert.NoError(t, err)

			// Verify the results
			assert.Equal(t, tt.expected, updateInfos)
		})
	}
}

func TestUpdateCheckerCheck(t *testing.T) {
	tests := []struct {
		name     string
		fileData string
		expected []UpdateInfo
	}{
		{
			name: "Single image",
			fileData: `
image: library/myimage:1.19.0
`,

			expected: []UpdateInfo{
				{
					RawLine:       "image: library/myimage:1.19.0",
					FullImageName: "library/myimage:1.19.0",
					ImageName:     "library/myimage",
					CurrentTag:    "1.19.0",
					LatestTag:     "1.20.0",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file with the test data
			file, err := os.CreateTemp("", "testfile.yaml")
			assert.NoError(t, err)
			defer os.Remove(file.Name())

			_, err = file.WriteString(tt.fileData)
			assert.NoError(t, err)
			file.Close()

			// Update the expected FilePath to match the temporary file name
			for i := range tt.expected {
				tt.expected[i].FilePath = file.Name()
			}

			// Create an UpdateChecker instance
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"count": 2, "results": [
					{"name": "1.18.0"},
					{"name": "1.18.1"},
					{"name": "1.19.0"},
					{"name": "1.20.0"}
				],"next": null}`))
			}))
			defer server.Close()

			registry := NewRegistry(server.URL)
			updateChecker := NewUpdateChecker(file.Name(), registry)

			// Call Check
			result, err := updateChecker.Check(true, true, true)
			assert.NoError(t, err)

			// Verify the results
			assert.Equal(t, tt.expected, result)
		})
	}
}
