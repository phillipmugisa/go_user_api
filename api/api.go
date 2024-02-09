package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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

func (s *ApiServer) Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// handle user creation
	r.Post("/user/create", makeHttpHandler(s.handleCreateUser))
	r.Post("/user/checkOtpcode", makeHttpHandler(s.handleUserVerification))

	fmt.Printf("Starting server on port: %s.\n", s.listenAddr)

	server := &http.Server{
		Addr:        s.listenAddr,
		Handler:     r,
		IdleTimeout: 120 * time.Second,
		ReadTimeout: 1 * time.Second,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// creating signal
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	fmt.Println("Recieved terminate, graceful shutdown", sig)

	// graceful shutdown
	tc, _ := context.WithDeadline(context.Background(), <-time.After(2*time.Second))
	server.Shutdown(tc)

}
