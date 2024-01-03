package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tkdn/go-study/log"
	"github.com/tkdn/go-study/middleware"
)

type JsonResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Query   int    `json:"query,omitempty"`
}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/", handler)

	if err := http.ListenAndServe(":8080", middleware.RejectAdminInternal(r)); err != nil {
		log.Logger.Error(err.Error())
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
	log.Logger.Info("root handler", "query", qs)
	if err := json.NewEncoder(w).Encode(s); err != nil {
		log.Logger.Error(err.Error())
	}
}

func notFoundHanlder(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	log.Logger.Info("not found handler")
	if _, err := w.Write([]byte("Not Found.")); err != nil {
		log.Logger.Error(err.Error())
	}
}
