package domain

type Storage interface {
    // Save a URL mapping. If the same Original already exists, returns the existing shortenedURL.
    Save(u *URL) (string, error)
    // Get original URL by shortcode
    Get(shortenedUrl string) (string, error)
    // Increment count for the domain_name for metrics
    IncDomainCount(domain_name string)
    // Get top N domains by count
    TopDomains(n int) map[string]int
}
