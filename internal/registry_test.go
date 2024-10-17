package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchImageTags_ValidImageWithTags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"count": 2, "results": [{"name": "latest"}, {"name": "18.04"}, {"name": "20.04"}, {"name": "22.04"}],"next": null}`))
	}))
	defer server.Close()

	registry := NewRegistry(server.URL)
	gotTags, err := registry.FetchImageTags("library/ubuntu")

	assert.NoError(t, err)
	assert.Equal(t, []string{"latest", "18.04","20.04","22.04"}, gotTags)
}
