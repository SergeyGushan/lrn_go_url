package urlhandlers

import (
	"github.com/SergeyGushan/lrn_go_url/internal/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGet(t *testing.T) {
	key := "test_key"
	value := "test_value"

	storage.URLStore.Push(key, value)

	_, exists := storage.URLStore.GetByKey(key)
	assert.True(t, exists)
}

func TestSave(t *testing.T) {
	key := "test_key"
	value := "test_value"

	storage.URLStore.Push(key, value)

	_, exists := storage.URLStore.GetByKey(key)
	assert.True(t, exists)
}
