package main

import (
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
	assert.Equal(t, urlStore[shortURL], dataValue)
}
func Test_getUrl(t *testing.T) {
	shortUrl := "/MeQpwyse"
	dataValue := "https://github.com/SergeyGushan"
	urlStore[host+shortUrl] = dataValue

	ts := httptest.NewServer(URLRouter())
	defer ts.Close()

	response, _ := testRequest(t, ts, http.MethodGet, shortUrl, "", true)

	assert.Equal(t, response.StatusCode, http.StatusTemporaryRedirect)
	assert.Equal(t, response.Header.Get("Location"), dataValue)
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
