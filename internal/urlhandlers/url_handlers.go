package urlhandlers

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/SergeyGushan/lrn_go_url/cmd/config"
	"github.com/SergeyGushan/lrn_go_url/internal/database"
	"github.com/SergeyGushan/lrn_go_url/internal/logger"
	"github.com/SergeyGushan/lrn_go_url/internal/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strings"
)

type StructReq struct {
	URL string `json:"url"`
}

type StructRes struct {
	Result string `json:"result"`
}

func PingDB(res http.ResponseWriter, req *http.Request) {
	err := database.DBClient.Ping()

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func Save(res http.ResponseWriter, req *http.Request) {
	var bodyReader io.Reader = req.Body

	if strings.Contains(req.Header.Get("Content-Encoding"), "gzip") {
		reader, err := gzip.NewReader(req.Body)
		if err != nil {
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer func(reader *gzip.Reader) {
			err := reader.Close()
			if err != nil {
				panic(err)
			}
		}(reader)
		bodyReader = reader
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// Обрабатываем ошибку
			panic(err)
		}
	}(req.Body)

	body, err := io.ReadAll(bodyReader)
	if err != nil {
		panic(err)
	}
	longURL := string(body)
	logger.Log.Info(longURL)
	hash := md5.New()

	_, err = io.WriteString(hash, longURL)
	if err != nil {
		return
	}

	shortCode := base64.URLEncoding.EncodeToString(hash.Sum(nil))[:8]

	shortURL := fmt.Sprintf("%s/%s", config.Opt.BaseURL, shortCode)

	storage.URLStore.Push(shortURL, longURL)

	res.WriteHeader(http.StatusCreated)

	_, err = res.Write([]byte(shortURL))
	if err != nil {
		return
	}
}

func Shorten(res http.ResponseWriter, req *http.Request) {
	var structReq StructReq
	var structRes StructRes
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)

	if err != nil {
		return
	}

	// десериализуем JSON в Visitor
	if err = json.Unmarshal(buf.Bytes(), &structReq); err != nil {
		return
	}

	longURL := structReq.URL

	if longURL == "" {
		return
	}

	hash := sha1.New()

	_, err = io.WriteString(hash, longURL)
	if err != nil {
		return
	}

	shortCode := base64.URLEncoding.EncodeToString(hash.Sum(nil))[:8]
	shortURL := fmt.Sprintf("%s/%s", config.Opt.BaseURL, shortCode)

	storage.URLStore.Push(shortURL, longURL)

	structRes.Result = shortURL
	respJSON, err := json.Marshal(structRes)
	if err != nil {
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

	_, err = res.Write(respJSON)
	if err != nil {
		return
	}

}

func Get(res http.ResponseWriter, req *http.Request) {
	shortCode := chi.URLParam(req, "shortCode")

	shortURL := fmt.Sprintf("%s/%s", config.Opt.BaseURL, shortCode)

	url, hasURL := storage.URLStore.GetByKey(shortURL)

	if hasURL {
		res.Header().Set("Location", url)
		res.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	http.Error(res, "Bad Request", http.StatusBadRequest)
}
