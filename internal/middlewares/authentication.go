package middlewares

import (
	"errors"
	"github.com/SergeyGushan/lrn_go_url/internal/authentication"
	"net/http"
	"time"
)

const tokenKey = "token"

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		_, err := getTokenFromCookie(req)
		var TokenError *authentication.TokenError

		if errors.As(err, &TokenError) {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err != nil {
			_, errSetToken := setTokenToCookie(res)
			if errSetToken != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		next.ServeHTTP(res, req)
	})
}

func getTokenFromCookie(req *http.Request) (string, error) {
	tokenCookie, err := req.Cookie(tokenKey)

	if err != nil {
		return "", err
	}

	token := tokenCookie.Value
	UserID, errUserID := authentication.GetUserIDFromJWTString(token)

	if errUserID != nil {
		return "", errUserID
	}

	authentication.User().SetID(UserID)

	return token, nil
}

func setTokenToCookie(res http.ResponseWriter) (string, error) {
	token, err := authentication.BuildJWTString(authentication.User().GetID())

	if err == nil {
		http.SetCookie(res, &http.Cookie{
			Name:    tokenKey,
			Value:   token,
			Expires: time.Now().Add(authentication.TokenExp),
			Path:    "/",
		})
	}

	return token, err
}
