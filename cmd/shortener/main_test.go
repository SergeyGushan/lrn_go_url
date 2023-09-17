package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func Test_saveUrl(t *testing.T) {
	data := url.Values{}
	dataKey := "url"
	dataValue := "https://github.com/SergeyGushan"

	data.Set(dataKey, dataValue)
	bodyData := data.Encode()

	requestPost := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(bodyData))
	requestPost.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	wPost := httptest.NewRecorder()

	shortenerHandler(wPost, requestPost)
	resPost := wPost.Result()

	assert.Equal(t, resPost.StatusCode, http.StatusCreated)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		assert.NoError(t, err)
	}(resPost.Body)

	resBody, err := io.ReadAll(resPost.Body)
	assert.NoError(t, err)

	shortUrl := strings.Trim(string(resBody), "/")

	assert.Equal(t, urlStore[shortUrl], dataValue)
}
func Test_getUrl(t *testing.T) {
	shortUrl := "http://localhost:8080/MeQpwyse"
	dataValue := "https://github.com/SergeyGushan"
	urlStore[shortUrl] = dataValue

	requestGet := httptest.NewRequest(http.MethodGet, shortUrl, nil)

	wGet := httptest.NewRecorder()
	shortenerHandler(wGet, requestGet)

	resGet := wGet.Result()

	assert.Equal(t, resGet.StatusCode, http.StatusTemporaryRedirect)
	assert.Equal(t, resGet.Header.Get("Location"), dataValue)
}
