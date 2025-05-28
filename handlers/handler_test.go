package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/tajimyradov/quotes-api/models"
)

type mockStorage struct {
	quotes []models.Quote
}

func (m *mockStorage) Create(q models.Quote) models.Quote {
	q.ID = len(m.quotes) + 1
	m.quotes = append(m.quotes, q)
	return q
}

func (m *mockStorage) GetAll() []models.Quote {
	return m.quotes
}

func (m *mockStorage) GetByAuthor(author string) []models.Quote {
	var filtered []models.Quote
	for _, q := range m.quotes {
		if q.Author == author {
			filtered = append(filtered, q)
		}
	}
	return filtered
}

func (m *mockStorage) GetRandom() (models.Quote, error) {
	if len(m.quotes) == 0 {
		return models.Quote{}, errors.New("no quotes available")
	}
	return m.quotes[0], nil
}

func (m *mockStorage) Delete(id int) error {
	for i, q := range m.quotes {
		if q.ID == id {
			m.quotes = append(m.quotes[:i], m.quotes[i+1:]...)
			return nil
		}
	}
	return errors.New("quote not found")
}
func TestCreateQuote(t *testing.T) {
	store := &mockStorage{}
	handler := NewQuoteHandler(store)

	body := []byte(`{"author":"Test","quote":"Hello"}`)
	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateQuote(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.Quote
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "Test", resp.Author)
	assert.Equal(t, "Hello", resp.Text)
}

func TestCreateQuoteInvalidJSON(t *testing.T) {
	store := &mockStorage{}
	handler := NewQuoteHandler(store)

	body := []byte(`{invalid json}`)
	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateQuote(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateQuoteEmptyFields(t *testing.T) {
	store := &mockStorage{}
	handler := NewQuoteHandler(store)

	body := []byte(`{"author":"","quote":""}`)
	req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateQuote(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetQuotesByAuthorFound(t *testing.T) {
	store := &mockStorage{
		quotes: []models.Quote{
			{ID: 1, Author: "A", Text: "Q1"},
			{ID: 2, Author: "B", Text: "Q2"},
		},
	}
	handler := NewQuoteHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/quotes?author=A", nil)
	w := httptest.NewRecorder()

	handler.GetAllQuotes(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []models.Quote
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "A", resp[0].Author)
}

func TestGetQuotesByAuthorNotFound(t *testing.T) {
	store := &mockStorage{}
	handler := NewQuoteHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/quotes?author=NonExistent", nil)
	w := httptest.NewRecorder()

	handler.GetAllQuotes(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAllQuotes(t *testing.T) {
	store := &mockStorage{
		quotes: []models.Quote{
			{ID: 1, Author: "A", Text: "Q1"},
			{ID: 2, Author: "B", Text: "Q2"},
		},
	}
	handler := NewQuoteHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/quotes", nil)
	w := httptest.NewRecorder()

	handler.GetAllQuotes(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []models.Quote
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
}

func TestGetRandomQuoteEmpty(t *testing.T) {
	store := &mockStorage{}
	handler := NewQuoteHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/quotes/random", nil)
	w := httptest.NewRecorder()

	handler.GetRandomQuote(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetRandomQuoteSuccess(t *testing.T) {
	store := &mockStorage{
		quotes: []models.Quote{{ID: 1, Author: "X", Text: "Y"}},
	}
	handler := NewQuoteHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/quotes/random", nil)
	w := httptest.NewRecorder()

	handler.GetRandomQuote(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.Quote
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.ID)
	assert.Equal(t, "X", resp.Author)
}

func TestDeleteQuoteSuccess(t *testing.T) {
	store := &mockStorage{
		quotes: []models.Quote{{ID: 1, Author: "X", Text: "Y"}},
	}
	handler := NewQuoteHandler(store)

	req := httptest.NewRequest(http.MethodDelete, "/quotes/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.DeleteQuote(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteQuoteNotFound(t *testing.T) {
	store := &mockStorage{}
	handler := NewQuoteHandler(store)

	req := httptest.NewRequest(http.MethodDelete, "/quotes/99", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "99"})
	w := httptest.NewRecorder()

	handler.DeleteQuote(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteQuoteInvalidID(t *testing.T) {
	store := &mockStorage{}
	handler := NewQuoteHandler(store)

	req := httptest.NewRequest(http.MethodDelete, "/quotes/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	handler.DeleteQuote(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
