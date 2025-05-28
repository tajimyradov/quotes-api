package storage

import (
	"errors"
	"math/rand"
	"strings"
	"sync"

	"github.com/tajimyradov/quotes-api/models"
)

type QuoteStorage interface {
	Create(models.Quote) models.Quote
	GetAll() []models.Quote
	GetRandom() (models.Quote, error)
	GetByAuthor(author string) []models.Quote
	Delete(id int) error
}

type MemoryStorage struct {
	sync.Mutex
	quotes []models.Quote
	nextID int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		quotes: []models.Quote{},
		nextID: 1,
	}
}

func (s *MemoryStorage) Create(q models.Quote) models.Quote {
	s.Lock()
	defer s.Unlock()
	q.ID = s.nextID
	s.nextID++
	s.quotes = append(s.quotes, q)
	return q
}

func (s *MemoryStorage) GetAll() []models.Quote {
	s.Lock()
	defer s.Unlock()
	return append([]models.Quote(nil), s.quotes...)
}

func (s *MemoryStorage) GetRandom() (models.Quote, error) {
	s.Lock()
	defer s.Unlock()
	if len(s.quotes) == 0 {
		return models.Quote{}, errors.New("no quotes available")
	}
	return s.quotes[rand.Intn(len(s.quotes))], nil
}

func (s *MemoryStorage) GetByAuthor(author string) []models.Quote {
	s.Lock()
	defer s.Unlock()
	var filtered []models.Quote
	for _, q := range s.quotes {
		if strings.EqualFold(q.Author, author) {
			filtered = append(filtered, q)
		}
	}
	return filtered
}

func (s *MemoryStorage) Delete(id int) error {
	s.Lock()
	defer s.Unlock()
	for i, q := range s.quotes {
		if q.ID == id {
			s.quotes = append(s.quotes[:i], s.quotes[i+1:]...)
			return nil
		}
	}
	return errors.New("quote not found")
}
