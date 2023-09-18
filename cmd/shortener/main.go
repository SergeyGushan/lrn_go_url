package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/SergeyGushan/lrn_go_url/cmd/config"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

var host string

func main() {
	options := config.Flags()
	host = options.B
	err := http.ListenAndServe(options.A, UrlRouter())
	if err != nil {
		return
	}
}

var urlStore = make(map[string]string)

func UrlRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/", saveUrlHandler)
	r.Get("/{shortCode}", getUrlHandler)

	return r
}

func saveUrlHandler(res http.ResponseWriter, req *http.Request) {
	longURL := req.FormValue("url")

	hash := md5.New()

	_, err := io.WriteString(hash, longURL)
	if err != nil {
		return
	}

	shortCode := base64.URLEncoding.EncodeToString(hash.Sum(nil))[:8]

	shortURL := fmt.Sprintf("%s/%s", host, shortCode)

	urlStore[shortURL] = longURL

	res.WriteHeader(http.StatusCreated)

	_, err = res.Write([]byte(shortURL))
	if err != nil {
		return
	}
}

func getUrlHandler(res http.ResponseWriter, req *http.Request) {
	shortCode := chi.URLParam(req, "shortCode")

	shortURL := fmt.Sprintf("%s/%s", host, shortCode)

	url, hasUrl := urlStore[shortURL]

	if hasUrl {
		res.Header().Set("Location", url)
		res.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	http.Error(res, "Bad Request", http.StatusBadRequest)
}
