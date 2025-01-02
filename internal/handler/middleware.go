package handler

import (
	"context"
	"net/http"
)

type SessionKey string

func SessionIDToContext(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var sessionID string
		cookie, err := r.Cookie("session-id")

		if err == nil {
			sessionID = cookie.Value
		}

		ctx := context.WithValue(r.Context(), SessionKey("session-id"), sessionID)
		newRequest := r.WithContext(ctx)
		next.ServeHTTP(w, newRequest)
	}
	return http.HandlerFunc(fn)
}
