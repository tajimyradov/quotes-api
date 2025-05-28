package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tajimyradov/quotes-api/handlers/utils"
	"github.com/tajimyradov/quotes-api/models"
	"github.com/tajimyradov/quotes-api/storage"
)

type QuoteHandler struct {
	db storage.QuoteStorage
}

func NewQuoteHandler(db storage.QuoteStorage) *QuoteHandler {
	return &QuoteHandler{db: db}
}

// CreateQuote godoc
// @Summary Add a new quote
// @Description Add a new quote with author and text
// @Tags quotes
// @Accept json
// @Produce json
// @Param quote body models.Quote true "Quote to add"
// @Success 200 {object} models.Quote
// @Failure 400 {object} map[string]string
// @Router /quotes [post]
func (h *QuoteHandler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var q models.Quote
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		log.Printf("CreateQuote: failed to decode request: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid input")
		return
	}
	if q.Author == "" || q.Text == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Author and quote cannot be empty")
		return
	}

	created := h.db.Create(q)
	log.Printf("Created quote ID %d by '%s'", created.ID, created.Author)
	utils.RespondWithJSON(w, http.StatusOK, created)
}

// GetAllQuotes godoc
// @Summary Get all quotes or filter by author
// @Description Retrieve all quotes or filter by author query parameter
// @Tags quotes
// @Accept json
// @Produce json
// @Param author query string false "Author to filter quotes"
// @Success 200 {array} models.Quote
// @Failure 404 {object} map[string]string
// @Router /quotes [get]
func (h *QuoteHandler) GetAllQuotes(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")

	var quotes []models.Quote
	if author != "" {
		log.Printf("Filtering quotes by author: %s", author)
		quotes = h.db.GetByAuthor(author)
		if len(quotes) == 0 {
			log.Printf("No quotes found for author: %s", author)
			utils.RespondWithError(w, http.StatusNotFound, "No quotes found for author")
			return
		}
	} else {
		log.Println("Retrieving all quotes")
		quotes = h.db.GetAll()
	}

	utils.RespondWithJSON(w, http.StatusOK, quotes)
}

// GetRandomQuote godoc
// @Summary Get a random quote
// @Description Retrieve a random quote from the store
// @Tags quotes
// @Accept json
// @Produce json
// @Success 200 {object} models.Quote
// @Failure 404 {object} map[string]string
// @Router /quotes/random [get]
func (h *QuoteHandler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	q, err := h.db.GetRandom()
	if err != nil {
		log.Println("No quotes available for random selection")
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	log.Printf("Retrieved random quote ID %d", q.ID)
	utils.RespondWithJSON(w, http.StatusOK, q)
}

// DeleteQuote godoc
// @Summary Delete quote by ID
// @Description Delete a quote by its ID
// @Tags quotes
// @Param id path int true "Quote ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /quotes/{id} [delete]
func (h *QuoteHandler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("DeleteQuote: invalid ID format: %s", idStr)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	err = h.db.Delete(id)
	if err != nil {
		log.Printf("DeleteQuote: quote with ID %d not found", id)
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	log.Printf("Deleted quote with ID %d", id)
	w.WriteHeader(http.StatusNoContent)
}
