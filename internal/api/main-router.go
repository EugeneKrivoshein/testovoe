package api

import (
	"net/http"
	"testovoe/internal/handlers"

	"github.com/gorilla/mux"
)

func NewRouter(authHandler *handlers.AuthHandler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Сервер работает!"))
	}).Methods("GET")

	router.HandleFunc("/api/register", authHandler.RegisterHandler).Methods("POST")
	router.HandleFunc("/api/auth", authHandler.AuthHandler).Methods("POST")
	router.HandleFunc("/api/auth/{token}", authHandler.LogoutHandler).Methods("POST")

	return router
}
