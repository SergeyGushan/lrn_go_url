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
	options := config.SetOptions()
	host = options.B
	err := http.ListenAndServe(options.A, URLRouter())
	if err != nil {
		return
	}
}

var urlStore = make(map[string]string)

func URLRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/", saveURLHandler)
	r.Get("/{shortCode}", getURLHandler)

	return r
}

func saveURLHandler(res http.ResponseWriter, req *http.Request) {
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

func getURLHandler(res http.ResponseWriter, req *http.Request) {
	shortCode := chi.URLParam(req, "shortCode")

	shortURL := fmt.Sprintf("%s/%s", host, shortCode)

	url, hasURL := urlStore[shortURL]

	if hasURL {
		res.Header().Set("Location", url)
		res.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	http.Error(res, "Bad Request", http.StatusBadRequest)
}
