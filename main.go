package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tkdn/go-study/infra/database"
	"github.com/tkdn/go-study/log"
	"github.com/tkdn/go-study/middleware"
)

type JsonResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Query   int            `json:"query,omitempty"`
	User    *database.User `json:"user"`
}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/", handler)

	if err := http.ListenAndServe(":8080", middleware.RejectAdminInternal(r)); err != nil {
		log.Logger.Error(err.Error())
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDB()
	defer db.Close()

	if r.URL.Path != "/" {
		notFoundHanlder(w, r)
		return
	}

	qs := r.URL.Query().Get("query")
	qi, _ := strconv.Atoi(qs)
	users := database.NewUserDB(db)
	u, err := users.GetById(1)
	if err != nil {
		log.Logger.Error(err.Error())
	}

	s := JsonResponse{
		Status:  "success",
		Message: "root handler",
		Query:   qi,
		User:    u,
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
