package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/tkdn/go-study/infra/database"
	"github.com/tkdn/go-study/log"
	"github.com/tkdn/go-study/middleware"
)

type JsonResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Query   int            `json:"query,omitempty"`
	User    *database.User `json:"user,omitempty"`
}

type server struct {
	db *sqlx.DB
}

func main() {
	db := database.ConnectDB()
	defer db.Close()

	s := &server{db}
	r := http.NewServeMux()
	r.HandleFunc("/", s.handler)
	m := middleware.RejectAdminInternal(r)

	if err := http.ListenAndServe(":8080", m); err != nil {
		log.Logger.Error(err.Error())
	}
}

func (s *server) handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFoundHanlder(w, r)
		return
	}

	var u *database.User
	q := r.URL.Query().Get("query")
	qi, _ := strconv.Atoi(q)
	id := r.URL.Query().Get("user_id")
	ui, err := strconv.Atoi(id)

	if err == nil {
		users := database.NewUserRepository(s.db)
		u, _ = users.GetById(ui)
	}

	res := JsonResponse{
		Status:  "success",
		Message: "root handler",
		Query:   qi,
		User:    u,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Logger.Info("root handler", "query", q, "user", u)
	json.NewEncoder(w).Encode(res)
}

func notFoundHanlder(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	log.Logger.Info("not found handler")
	if _, err := w.Write([]byte("Not Found.")); err != nil {
		log.Logger.Error(err.Error())
	}
}
