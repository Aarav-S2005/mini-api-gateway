package middleware

import (
	"net/http"
)

func CookieToBearer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("x-auth-token")
		if err == nil && cookie.Value != "" {
			r.Header.Set("Authorization", "Bearer "+cookie.Value)
		}
		next.ServeHTTP(w, r)
	})
}
