package urlhandlers

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/SergeyGushan/lrn_go_url/cmd/config"
	"github.com/SergeyGushan/lrn_go_url/internal/storage"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func Save(res http.ResponseWriter, req *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(req.Body)
	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	longURL := string(body)

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
