package infrastructure

import (
	"testing"

	"errors"

	"github.com/MirMonajir/mir-url-shortener/internal_logic/domain"
	"github.com/stretchr/testify/assert"
)

func TestSave_NewURL_ReturnsShortenedUrl(t *testing.T) {
	store := NewInMemoryStore()

	urlStr := "https://example.com"
	urlObj, err := domain.NewURL(urlStr)
	assert.NoError(t, err)

	shortCode, err := store.Save(urlObj)

	assert.NoError(t, err)
	assert.NotEmpty(t, shortCode)
	assert.Equal(t, shortCode, urlObj.ShortenedUrl)
}

func TestSave_DuplicateURL_ReturnsSameShortenedUrl(t *testing.T) {
	store := NewInMemoryStore()

	urlStr := "https://example.com"
	urlObj1, _ := domain.NewURL(urlStr)
	urlObj2, _ := domain.NewURL(urlStr)

	short1, _ := store.Save(urlObj1)
	short2, _ := store.Save(urlObj2)

	assert.Equal(t, short1, short2)
}

func TestGet_ValidShortUrl_ReturnsOriginalUrl(t *testing.T) {
	store := NewInMemoryStore()

	urlStr := "https://example.com"
	urlObj, _ := domain.NewURL(urlStr)

	shortCode, _ := store.Save(urlObj)

	original, err := store.Get(shortCode)

	assert.NoError(t, err)
	assert.Equal(t, urlStr, original)
}

func TestGet_InvalidShortUrl_ReturnsError(t *testing.T) {
	store := NewInMemoryStore()

	original, err := store.Get("nonexistent")

	assert.Error(t, err)
	assert.Equal(t, "", original)
	assert.Equal(t, errors.New("shortened url not found").Error(), err.Error())
}

func TestIncDomainCount_ReturnsManuallyIncrements(t *testing.T) {
	store := NewInMemoryStore()

	domain := "example.com"
	store.IncDomainCount(domain)
	store.IncDomainCount(domain)

	top := store.TopDomains(1)

	assert.Equal(t, 2, top[domain])
}

func TestTopDomains_ReturnsSortedTopN(t *testing.T) {
	store := NewInMemoryStore()

	// Increment domain counts
	for i := 0; i < 3; i++ {
		store.IncDomainCount("a.com")
	}
	for i := 0; i < 5; i++ {
		store.IncDomainCount("b.com")
	}
	for i := 0; i < 2; i++ {
		store.IncDomainCount("c.com")
	}

	top := store.TopDomains(2)

	assert.Len(t, top, 2)
	assert.Equal(t, 5, top["b.com"])
	assert.Equal(t, 3, top["a.com"])
}

func TestSave_IncrementsDomainCount(t *testing.T) {
	store := NewInMemoryStore()

	urlStr := "https://sub.google.com/page"
	urlObj, _ := domain.NewURL(urlStr)

	_, err := store.Save(urlObj)
	assert.NoError(t, err)

	// "google.com" should be extracted and counted
	top := store.TopDomains(1)
	assert.Equal(t, 1, top["google.com"])
}
