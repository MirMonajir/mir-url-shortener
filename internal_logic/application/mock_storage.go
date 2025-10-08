package application

import (
	"github.com/MirMonajir/mir-url-shortener/internal_logic/domain"
	"github.com/stretchr/testify/mock"
)

// MockStorage is a mock implementation of the domain.Storage interface
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Save(u *domain.URL) (string, error) {
	args := m.Called(u)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) Get(code string) (string, error) {
	args := m.Called(code)
	return args.String(0), args.Error(1)
}

func (m *MockStorage) TopDomains(n int) map[string]int {
	args := m.Called(n)
	return args.Get(0).(map[string]int)
}

func (m *MockStorage) IncDomainCount(domain string) {
	m.Called(domain)
}
