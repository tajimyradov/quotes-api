package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRespondWithJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	payload := map[string]string{"message": "hello"}

	RespondWithJSON(rr, http.StatusOK, payload)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var response map[string]string
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "hello", response["message"])
}

func TestRespondWithError(t *testing.T) {
	rr := httptest.NewRecorder()

	RespondWithError(rr, http.StatusBadRequest, "something went wrong")

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var response map[string]string
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "something went wrong", response["error"])
}
