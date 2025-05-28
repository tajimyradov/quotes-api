package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tajimyradov/quotes-api/models"
)

func TestCreateAndGetAll(t *testing.T) {
	store := NewMemoryStorage()

	q := models.Quote{Author: "A", Text: "B"}
	store.Create(q)
	all := store.GetAll()

	assert.Equal(t, 1, len(all))
	assert.Equal(t, "A", all[0].Author)
	assert.Equal(t, "B", all[0].Text)
}

func TestGetByAuthor(t *testing.T) {
	store := NewMemoryStorage()
	store.Create(models.Quote{Author: "A", Text: "1"})
	store.Create(models.Quote{Author: "B", Text: "2"})

	result := store.GetByAuthor("A")

	assert.Len(t, result, 1)
	assert.Equal(t, "1", result[0].Text)
}

func TestGetRandom_Empty(t *testing.T) {
	store := NewMemoryStorage()

	_, err := store.GetRandom()
	assert.Error(t, err)
}

func TestDelete(t *testing.T) {
	store := NewMemoryStorage()
	store.Create(models.Quote{Author: "X", Text: "Y"})
	err := store.Delete(1)
	assert.NoError(t, err)

	err = store.Delete(1)
	assert.Error(t, err)
}
