package api

import (
	"net/http"
	"testovoe/internal/db"
	"testovoe/internal/handlers"
	"testovoe/internal/middleware"

	"github.com/gorilla/mux"
)

func NewRouter(authHandler *handlers.AuthHandler, docHandler *handlers.DocumentHandler, dbProvider *db.PostgresProvider) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Сервер работает!"))
	}).Methods("GET")

	router.HandleFunc("/api/register", authHandler.RegisterHandler).Methods("POST")
	router.HandleFunc("/api/auth", authHandler.AuthHandler).Methods("POST")
	router.HandleFunc("/api/auth/{token}", authHandler.LogoutHandler).Methods("POST")
	router.HandleFunc("/api/docs", docHandler.UploadDocumentHandler).Methods("POST").Handler(middleware.TokenAuthMiddleware(dbProvider)(http.HandlerFunc(docHandler.UploadDocumentHandler)))

	return router
}
