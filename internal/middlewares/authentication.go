package middlewares

import (
	"context"
	"errors"
	"github.com/SergeyGushan/lrn_go_url/internal/authentication"
	"net/http"
	"time"
)

const TokenKey = "token"

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		user := authentication.User{}
		_, err := getTokenFromCookie(req, &user)
		var TokenError *authentication.TokenError

		if errors.As(err, &TokenError) {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		if err != nil {
			_, errSetToken := setTokenToCookie(res, &user)
			if errSetToken != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		ctx := context.WithValue(req.Context(), "userID", user.ID)
		req = req.WithContext(ctx)
		next.ServeHTTP(res, req)
	})
}

func getTokenFromCookie(req *http.Request, user *authentication.User) (string, error) {
	tokenCookie, err := req.Cookie(TokenKey)
	if err != nil {
		return "", err
	}

	token := tokenCookie.Value
	UserID, errUserID := authentication.GetUserIDFromJWTString(token)

	if errUserID != nil {
		return "", errUserID
	}

	user.SetID(UserID)

	return token, nil
}

func setTokenToCookie(res http.ResponseWriter, user *authentication.User) (string, error) {
	token, err := authentication.BuildJWTString(user.GetID())

	if err == nil {
		http.SetCookie(res, &http.Cookie{
			Name:    TokenKey,
			Value:   token,
			Expires: time.Now().Add(authentication.TokenExp),
			Path:    "/",
		})
	}

	return token, err
}
