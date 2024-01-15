package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SergeyGushan/lrn_go_url/cmd/config"
	"github.com/SergeyGushan/lrn_go_url/internal/database"
	"github.com/SergeyGushan/lrn_go_url/internal/logger"
	"github.com/SergeyGushan/lrn_go_url/internal/storage"
	"github.com/SergeyGushan/lrn_go_url/internal/url"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func UserUrls(res http.ResponseWriter, req *http.Request) {
	userIDValue := req.Context().Value("userID")

	if userIDValue != nil {
		userID, ok := userIDValue.(string)
		if !ok {
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		urls := storage.Service.GetURLByUserID(userID)

		if len(urls) == 0 {
			http.Error(res, "No Content", http.StatusNoContent)
			return
		}

		respJSON, err := json.Marshal(urls)
		if err != nil {
			logger.Log.Error(err.Error())
			http.Error(res, "Bad Request", http.StatusBadRequest)
			return
		}

		res.Header().Set("Content-Type", "application/json")

		_, err = res.Write(respJSON)
		if err != nil {
			logger.Log.Error(err.Error())
			http.Error(res, "Bad Request", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}
}

func PingDB(res http.ResponseWriter, req *http.Request) {
	err := database.Client().Ping()

	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

func Save(res http.ResponseWriter, req *http.Request) {
	var bodyReader io.Reader = req.Body

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
	shortURL, errBuildShortURL := url.CreateShortURL(longURL)
	if errBuildShortURL != nil {
		return
	}

	userIDValue := req.Context().Value("userID")
	if userIDValue != nil {
		userID, ok := userIDValue.(string)
		if !ok {
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		errStorageSave := storage.Service.Save(shortURL, longURL, userID)
		var duplicateError *storage.DuplicateError
		if errors.As(errStorageSave, &duplicateError) {
			res.WriteHeader(http.StatusConflict)
			_, err = res.Write([]byte(duplicateError.ShortURL))
			if err != nil {
				return
			}
			return
		}

		if errStorageSave != nil {
			logger.Log.Error(errStorageSave.Error())
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusCreated)

		_, err = res.Write([]byte(shortURL))
		if err != nil {
			return
		}
	} else {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
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

	shortURL, errBuildShortURL := url.CreateShortURL(longURL)

	if errBuildShortURL != nil {
		return
	}

	userIDValue := req.Context().Value("userID")
	if userIDValue != nil {
		userID, ok := userIDValue.(string)
		if !ok {
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		errStorageSave := storage.Service.Save(shortURL, longURL, userID)

		var duplicateError *storage.DuplicateError
		if errors.As(errStorageSave, &duplicateError) {
			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(http.StatusConflict)
			structRes.Result = duplicateError.ShortURL
			respJSON, err := json.Marshal(structRes)
			if err != nil {
				return
			}
			_, err = res.Write(respJSON)
			if err != nil {
				return
			}
			return
		}

		if errStorageSave != nil {
			logger.Log.Error(errStorageSave.Error())
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

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
	} else {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}
}

func Get(res http.ResponseWriter, req *http.Request) {
	shortCode := chi.URLParam(req, "shortCode")

	shortURL := fmt.Sprintf("%s/%s", config.Opt.BaseURL, shortCode)

	URL, errStorageGet := storage.Service.GetOriginalURL(shortURL)
	if errStorageGet != nil {
		logger.Log.Error(errStorageGet.Error())
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", URL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func BatchCreate(res http.ResponseWriter, req *http.Request) {
	var batchReq []storage.BatchItemReq
	var batch []storage.BatchItem
	err := json.NewDecoder(req.Body).Decode(&batchReq)
	if err != nil {
		logger.Log.Error(err.Error())
		http.Error(res, "Bad Request", http.StatusBadRequest)
		return
	}

	userIDValue := req.Context().Value("userID")
	if userIDValue != nil {
		userID, ok := userIDValue.(string)
		if !ok {
			http.Error(res, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		for _, item := range batchReq {
			shortURL, errBuildShortURL := url.CreateShortURL(item.OriginalURL)
			if errBuildShortURL != nil {
				continue
			}

			batch = append(batch, storage.BatchItem{
				UserID:        userID,
				CorrelationID: item.CorrelationID,
				OriginalURL:   item.OriginalURL,
				ShortURL:      shortURL,
			})
		}

		results, err := storage.Service.SaveBatch(batch)

		if err != nil {
			logger.Log.Error(err.Error())
			http.Error(res, "Bad Request", http.StatusBadRequest)
			return
		}

		respJSON, err := json.Marshal(results)
		if err != nil {
			logger.Log.Error(err.Error())
			http.Error(res, "Bad Request", http.StatusBadRequest)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)

		_, err = res.Write(respJSON)
		if err != nil {
			logger.Log.Error(err.Error())
			http.Error(res, "Bad Request", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(res, "Unauthorized", http.StatusUnauthorized)
		return
	}
}
