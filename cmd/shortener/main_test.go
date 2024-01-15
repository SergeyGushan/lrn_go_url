package main

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/SergeyGushan/lrn_go_url/cmd/config"
	"github.com/SergeyGushan/lrn_go_url/internal/authentication"
	"github.com/SergeyGushan/lrn_go_url/internal/handlers"
	"github.com/SergeyGushan/lrn_go_url/internal/middlewares"
	"github.com/SergeyGushan/lrn_go_url/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

var fileName = os.TempDir() + "/test.log"

func Test_saveUrl(t *testing.T) {
	dataValue := "https://github.com/SergeyGushan"
	storage.Service, _ = storage.NewJSONStorage(fileName)

	ts := httptest.NewServer(URLRouter())
	defer ts.Close()

	requestPost, shortURL := testRequest(t, ts, http.MethodPost, "/", dataValue, false)

	assert.Equal(t, requestPost.StatusCode, http.StatusCreated)
	fullURL, err := storage.Service.GetOriginalURL(shortURL)
	assert.NoError(t, err)
	assert.Equal(t, fullURL, dataValue)

	defer func() {
		err := requestPost.Body.Close()
		require.NoError(t, err)
	}()
}

func Test_shortUrl(t *testing.T) {
	structRes := handlers.StructReq{}
	storage.Service, _ = storage.NewJSONStorage(fileName)
	structRes.URL = "https://github.com/SergeyGushan"
	respJSON, err := json.Marshal(structRes)
	if err != nil {
		return
	}

	ts := httptest.NewServer(URLRouter())
	defer ts.Close()

	requestPost, shortURL := testRequest(t, ts, http.MethodPost, "/api/shorten", string(respJSON), false)

	assert.Equal(t, requestPost.StatusCode, http.StatusCreated)
	fullURL, err := storage.Service.GetOriginalURL(shortURL)
	assert.NoError(t, err)
	assert.Equal(t, fullURL, structRes.URL)

	defer func() {
		err := requestPost.Body.Close()
		require.NoError(t, err)
	}()
}

func Test_getUrl(t *testing.T) {
	shortURL := "/MeQpwyse"
	dataValue := "https://github.com/SergeyGushan"
	storage.Service, _ = storage.NewJSONStorage(fileName)
	err := storage.Service.Save(config.Opt.BaseURL+shortURL, dataValue, "")
	assert.NoError(t, err)

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
	var structReq handlers.StructRes

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

	token, err := authentication.BuildJWTString(uuid.New().String())

	if err == nil {
		req.AddCookie(&http.Cookie{
			Name:    middlewares.TokenKey,
			Value:   token,
			Expires: time.Now().Add(authentication.TokenExp),
			Path:    "/",
		})
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	if resp.Header.Get("Content-Type") == "application/json" {
		if err = json.Unmarshal(respBody, &structReq); err != nil {
			panic(err)
		}

		respBody = []byte(structReq.Result)
	}

	return resp, string(respBody)
}

func TestGetOriginalURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ds := storage.NewDatabaseStorage(db)
	rows := sqlmock.NewRows([]string{"original_url", "is_deleted"}).AddRow("http://example.com", false)

	mock.ExpectQuery("SELECT original_url, is_deleted FROM urls WHERE short_url = \\$1").
		WithArgs("testShortURL").
		WillReturnRows(rows)

	originalURL, err := ds.GetOriginalURL("testShortURL")

	assert.NoError(t, err)
	assert.Equal(t, "http://example.com", originalURL)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSave(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ds := storage.NewDatabaseStorage(db)
	mock.ExpectExec("INSERT INTO urls \\(user_id, short_url, original_url\\) VALUES \\(\\$1, \\$2, \\$3\\) ON CONFLICT \\(original_url\\) DO NOTHING").
		WithArgs("", "testShortURL", "http://example.com").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = ds.Save("testShortURL", "http://example.com", "")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveBatch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ds := storage.NewDatabaseStorage(db)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO urls \\(user_id, correlation_id, short_url, original_url\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)").
		WithArgs("", "correlation1", "testShortURL1", "http://example.com1").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("INSERT INTO urls \\(user_id, correlation_id, short_url, original_url\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)").
		WithArgs("", "correlation2", "testShortURL2", "http://example.com2").
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	batch := []storage.BatchItem{
		{UserID: "", CorrelationID: "correlation1", ShortURL: "testShortURL1", OriginalURL: "http://example.com1"},
		{UserID: "", CorrelationID: "correlation2", ShortURL: "testShortURL2", OriginalURL: "http://example.com2"},
	}

	results, err := ds.SaveBatch(batch)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(results))
	assert.Equal(t, "correlation1", results[0].CorrelationID)
	assert.Equal(t, "/testShortURL1", results[0].ShortURL)
	assert.Equal(t, "correlation2", results[1].CorrelationID)
	assert.Equal(t, "/testShortURL2", results[1].ShortURL)

	assert.NoError(t, mock.ExpectationsWereMet())
}
