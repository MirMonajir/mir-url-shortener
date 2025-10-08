package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortenAndRedirect(t *testing.T) {
	r := setupRouter()

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
	assert.Contains(t, resp.ShortURL, "https://mir.com/")

	code := strings.TrimPrefix(resp.ShortURL, "https://mir.com/")

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
