package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenAndRedirect(t *testing.T) {
	// Ensure BASE_URL is set
	os.Setenv("SERVER_URL", "localhost:8080")
	defer os.Unsetenv("SERVER_URL")
	r := SetupRouter()

	// Prepare the json body for the POST request
	jsonBody := `{"URL":"https://google.com"}`
	// hit the post endpoint
	req := httptest.NewRequest("POST", "/shortenurl", bytes.NewBufferString(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	type ShortenResponse struct {
		ShortURL string `json:"short_url"`
	}

	var resp ShortenResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp.ShortURL, "http://localhost:8080/")

	code := strings.TrimPrefix(resp.ShortURL, "http://localhost:8080/")

	// Redirect request
	req = httptest.NewRequest("GET", "/"+code, nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Contains(t, []int{http.StatusMovedPermanently, http.StatusFound}, w.Code)
	assert.Equal(t, "https://google.com", w.Header().Get("Location"))

	// Metrics request
	req = httptest.NewRequest("GET", "/appmetrics", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "google.com")
}
