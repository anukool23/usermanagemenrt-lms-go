package middleware

import (
	"fmt"
	"net/http"

	"github.com/anukool23/usermanagement-lms-go/internal/utils/response"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func SecretKeyAuth(allowedKeys []string) Middleware {
	allowed := make(map[string]struct{}, len(allowedKeys))
	for _, key := range allowedKeys {
		allowed[key] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			secretKey := r.Header.Get("x-secret-key")
			if _, ok := allowed[secretKey]; !ok {
				_ = response.WriteJSON(w, http.StatusUnauthorized, response.GeneralError(fmt.Errorf("unauthorized")))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
