package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
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

type handler struct {
	db *sqlx.DB
}

var _ http.Handler = (*handler)(nil)

func main() {
	db := database.ConnectDB()
	defer db.Close()

	h := &handler{db}
	server := &http.Server{
		Addr:    net.JoinHostPort("", "8080"),
		Handler: middleware.RejectAdminInternal(h),
	}
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Logger.Error(err.Error())
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		users := database.NewUserRepository(h.db)
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
	fmt.Fprint(w, "Not Found.")
}
