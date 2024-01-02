package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/tkdn/go-study/middleware"
)

type JsonResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Query   int    `json:"query,omitempty"`
}

var opt = slog.HandlerOptions{}
var logger = slog.New(slog.NewJSONHandler(os.Stdout, &opt))

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/", handler)

	if err := http.ListenAndServe(":8080", middleware.HostCheckMiddleware(r)); err != nil {
		logger.Error(err.Error())
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFoundHanlder(w, r)
		return
	}
	qs := r.URL.Query().Get("query")
	i, _ := strconv.Atoi(qs)
	s := JsonResponse{
		Status:  "success",
		Message: "root handler",
		Query:   i,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	logger.Info("root handler", "query", qs)
	if err := json.NewEncoder(w).Encode(s); err != nil {
		logger.Error(err.Error())
	}
}

func notFoundHanlder(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	logger.Info("not found handler")
	if _, err := w.Write([]byte("Not Found.")); err != nil {
		logger.Error(err.Error())
	}
}
