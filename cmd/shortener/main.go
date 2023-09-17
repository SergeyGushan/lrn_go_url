package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

const host = "http://localhost:8080"

var urlStore = make(map[string]string)

func main() {
	run()
}

func run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", shortenerHandler)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return
	}
}

func saveUrl(res http.ResponseWriter, req *http.Request) {
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

func getUrl(res http.ResponseWriter, req *http.Request) {
	shortURL := fmt.Sprintf("%s%s", host, req.URL.RequestURI())

	url, hasUrl := urlStore[shortURL]

	if hasUrl {
		res.Header().Set("Location", url)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func shortenerHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		saveUrl(res, req)
		return
	}

	if req.Method == http.MethodGet {
		getUrl(res, req)
		return
	}

	http.Error(res, "Bad Request", http.StatusBadRequest)
}
