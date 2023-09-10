package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func shortenerHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		longURL := req.FormValue("url")
		hash := md5.New()

		_, err := io.WriteString(hash, longURL)
		if err != nil {
			return
		}

		shortCode := base64.URLEncoding.EncodeToString(hash.Sum(nil))[:8]

		urlStore[shortCode] = longURL
		shortURL := fmt.Sprintf("%s/%s", host, shortCode)

		_, err = res.Write([]byte(shortURL))
		if err != nil {
			return
		}

		res.WriteHeader(http.StatusCreated)
		return
	}

	if req.Method == http.MethodGet {
		shortCode := strings.Trim(req.RequestURI, "/")
		shortURL, hasUrl := urlStore[shortCode]

		if hasUrl {
			res.Header().Set("Location", shortURL)
			res.WriteHeader(http.StatusTemporaryRedirect)

			return
		}
	}

	http.Error(res, "Bad Request", http.StatusBadRequest)
}
