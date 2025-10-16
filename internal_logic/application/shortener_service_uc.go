package application

import (
	"fmt"
	"os"

	"github.com/MirMonajir/mir-url-shortener/internal_logic/domain"
)

type ShortenerService struct {
	repo domain.Storage
}

func NewShortenerService(r domain.Storage) *ShortenerService {
	return &ShortenerService{repo: r}
}

// Shorten the original URL
func (s *ShortenerService) Shorten(original string) (string, error) {
	u, err := domain.NewURL(original)
	if err != nil {
		return "", err
	}
	code, err := s.repo.Save(u)
	if err != nil {
		return "", err
	}
	full_url := os.Getenv("SERVER_URL")
	// return full URL, e.g. https://mir.com/{code}
	return fmt.Sprintf("http://%s/%s", full_url, code), nil
}

func (s *ShortenerService) Resolve(shortUrl string) (string, error) {
	return s.repo.Get(shortUrl)
}

func (s *ShortenerService) TopDomains(n int) map[string]int {
	return s.repo.TopDomains(n)
}
