package urlhandlers

import (
	"github.com/SergeyGushan/lrn_go_url/internal/storage"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var fileName = os.TempDir() + "/test.log"

func TestGet(t *testing.T) {
	key := "test_key"
	value := "test_value"
	storage.URLStore, _ = storage.NewURL(fileName)
	storage.URLStore.Push(key, value)

	_, exists := storage.URLStore.GetByKey(key)
	assert.True(t, exists)
}

func TestSave(t *testing.T) {
	key := "test_key"
	value := "test_value"
	storage.URLStore, _ = storage.NewURL(fileName)
	storage.URLStore.Push(key, value)

	_, exists := storage.URLStore.GetByKey(key)
	assert.True(t, exists)
}
