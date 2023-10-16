package main

import (
	"github.com/SergeyGushan/lrn_go_url/cmd/config"
	"github.com/SergeyGushan/lrn_go_url/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	ts := httptest.NewServer(URLRouter())
	defer ts.Close()

	requestPost, shortURL := testRequest(t, ts, http.MethodPost, "/", bodyData, false)

	assert.Equal(t, requestPost.StatusCode, http.StatusCreated)
	fullURL, hasURL := storage.URLStore.GetByKey(shortURL)
	assert.Equal(t, hasURL, true)
	assert.Equal(t, fullURL, dataValue)

	defer func() {
		err := requestPost.Body.Close()
		require.NoError(t, err)
	}()
}
func Test_getUrl(t *testing.T) {
	shortURL := "/MeQpwyse"
	dataValue := "https://github.com/SergeyGushan"
	storage.URLStore.Push(config.Opt.BaseUrl+shortURL, dataValue)

	ts := httptest.NewServer(URLRouter())
	defer ts.Close()

	response, _ := testRequest(t, ts, http.MethodGet, shortURL, "", true)

	assert.Equal(t, response.StatusCode, http.StatusTemporaryRedirect)
	assert.Equal(t, response.Header.Get("Location"), dataValue)

	defer func() {
		err := response.Body.Close()
		require.NoError(t, err)
	}()
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string, prohibitRedirects bool) (*http.Response, string) {

	var req *http.Request
	var err error

	if len(body) > 0 {
		req, err = http.NewRequest(method, ts.URL+path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req, err = http.NewRequest(method, ts.URL+path, nil)
	}

	require.NoError(t, err)

	if prohibitRedirects {
		ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Запретить редиректы
		}
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
