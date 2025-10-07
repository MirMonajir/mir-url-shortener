package infrastructure

import (
	"errors"
	"sync"

	"github.com/MirMonajir/mir-url-shortener/internal/domain"
)

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
	// generate a new shortened URL
	shortUrl := generateShortUrl(len(s.originalToShort) + 1)
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

func generateShortUrl(id int) string {
	// convert id to base62 string etc.
	return fmt.Sprintf("%d", id)
}

func extractDomain(original string) string {
	u, err := url.Parse(original)
	if err != nil {
		return ""
	}
	return u.Hostname()
}
