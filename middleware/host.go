package middleware

import (
	"fmt"
	"net/http"
)

func HostCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("my first middleware")
		next.ServeHTTP(w, r)
	})
}
