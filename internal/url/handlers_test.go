package url

import (
	"github.com/SergeyGushan/lrn_go_url/internal/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGet(t *testing.T) {
	key := "test_key"
	value := "test_value"
	storage.Service, _ = storage.NewStorage()
	err := storage.Service.Save(key, value, "")
	assert.NoError(t, err)

	_, exists := storage.Service.GetOriginalURL(key)
	assert.NoError(t, exists)
}

func TestSave(t *testing.T) {
	key := "test_key"
	value := "test_value"
	storage.Service, _ = storage.NewStorage()

	err := storage.Service.Save(key, value, "")
	assert.NoError(t, err)

	_, exists := storage.Service.GetOriginalURL(key)
	assert.NoError(t, exists)
}
