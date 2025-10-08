package domain

type Shortener interface {
	Shorten(originalUrl string) (string, error)
	Resolve(shortenedUrl string) (string, error)
	TopDomains(n int) map[string]int
}
