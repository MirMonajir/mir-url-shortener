package infrastructure

import (
	"errors"
	"math/rand"
	"net/url"
	"sort"
	"sync"

	"github.com/MirMonajir/mir-url-shortener/internal_logic/domain"
	"golang.org/x/net/publicsuffix"
)

const characterset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type InMemoryStore struct {
	mutex           sync.RWMutex
	originalToShort map[string]string
	shortToOriginal map[string]string
	domainCounts    map[string]int
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		originalToShort: make(map[string]string),
		shortToOriginal: make(map[string]string),
		domainCounts:    make(map[string]int),
	}
}

// Save implements storage
func (s *InMemoryStore) Save(u *domain.URL) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if short, exists := s.originalToShort[u.OriginalUrl]; exists {
		return short, nil
	}
	var shortUrl string
	for {
		// generate a new shortened url
		shortUrl = generateShortUrl()

		// Ensure the shortUrl hasn't already been used
		if _, exists := s.shortToOriginal[shortUrl]; !exists {
			break // It's unique, safe to use
		}
		// Otherwise, loop and generate a new one
	}
	u.ShortenedUrl = shortUrl
	s.originalToShort[u.OriginalUrl] = shortUrl
	s.shortToOriginal[shortUrl] = u.OriginalUrl

	// count the domain's
	domain := extractDomain(u.OriginalUrl)
	s.domainCounts[domain]++

	return shortUrl, nil
}

// Get implements Storage
func (s *InMemoryStore) Get(shortenedUrl string) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	original, exists := s.shortToOriginal[shortenedUrl]
	if !exists {
		return "", errors.New("shortened url not found")
	}
	return original, nil
}

// Increments the domain name count if already present
func (s *InMemoryStore) IncDomainCount(domain_name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.domainCounts[domain_name]++
}

// Returns the shortened domains name and count for metrics
func (s *InMemoryStore) TopDomains(n int) map[string]int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	//sort by counts
	result := make(map[string]int)

	type keyValue struct {
		Key   string
		Value int
	}
	var slices []keyValue
	for k, v := range s.domainCounts {
		slices = append(slices, keyValue{k, v})
	}
	// sort in descending order
	sort.Slice(slices, func(i, j int) bool {
		return slices[i].Value > slices[j].Value
	})
	for i, item := range slices {
		if i >= n {
			break
		}
		result[item.Key] = item.Value
	}
	return result
}

// helper functions

// generates the shortenedCode
func generateShortUrl() string {
	length := 6
	b := make([]byte, length)
	for i := range b {
		b[i] = characterset[rand.Intn(len(characterset))]
	}
	return string(b)
}

func extractDomain(original string) string {
	u, err := url.Parse(original)
	if err != nil {
		return ""
	}

	host := u.Hostname()
	domain, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		// fallback to original host if error
		return host
	}
	return domain
}
