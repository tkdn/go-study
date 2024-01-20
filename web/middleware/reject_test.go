package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tkdn/go-study/web/middleware"
)

var tests = []struct {
	name string
	host string
	want int
}{
	{"Forbidden from Host: 'admin.internal' header", "admin.internal", http.StatusForbidden},
	{"Allowed from Host: 'public.example' header", "public.example", http.StatusOK},
}

func TestRejectAdminInternal(t *testing.T) {
	for _, tt := range tests {
		r := http.NewServeMux()
		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte("Hello, world"))
			if err != nil {
				t.Errorf(err.Error())
			}
		})
		ts := httptest.NewServer(middleware.RejectAdminInternal(r))
		t.Cleanup(func() {
			ts.Close()
		})

		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", ts.URL, nil)
			if err != nil {
				t.Errorf(err.Error())
			}
			req.Host = tt.host
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Errorf(err.Error())
			}
			got := res.StatusCode
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
