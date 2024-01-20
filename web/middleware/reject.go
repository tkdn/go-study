package middleware

import (
	"net/http"

	"github.com/tkdn/go-study/log"
)

func RejectAdminInternal(next http.Handler) http.Handler {
	const rejectHost = "admin.internal"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host == rejectHost {
			w.WriteHeader(http.StatusForbidden)
			if _, err := w.Write([]byte("Fobidden Request")); err != nil {
				log.Logger.Error(err.Error())
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}
