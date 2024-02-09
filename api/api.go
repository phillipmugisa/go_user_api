package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/phillipmugisa/go_user_api/storage"
)

type ApiServer struct {
	listenAddr string
	store      storage.Storage
}

func NewApiServer(listenAddr string, store storage.Storage) *ApiServer {
	return &ApiServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *ApiServer) Run() error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// handle user creation
	r.Post("/user/create", makeHttpHandler(s.handleCreateUser))
	r.Post("/user/checkOtpcode", makeHttpHandler(s.handleUserVerification))

	fmt.Printf("Starting server on port: %s.\n", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, r)

}
