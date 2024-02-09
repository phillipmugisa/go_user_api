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
	// we are use go chi-chi as our router
	r := chi.NewRouter()

	// for purposes of logging
	r.Use(middleware.Logger)

	// api routes
	// handle user creation
	r.Post("/user/create", makeHttpHandler(s.handleCreateUser))
	// handle account verification using code sent to emain
	r.Post("/user/checkOtpcode", makeHttpHandler(s.handleUserVerification))

	server := &http.Server{
		Addr:        s.listenAddr,
		Handler:     r,
		IdleTimeout: 120 * time.Second,
		ReadTimeout: 5 * time.Second,
	}
	// WriteTimeout left out due to slow network during production,
	// while sending verification code

	// starting our server in a go routine
	go func() {
		fmt.Printf("Starting server on port: %s\n", s.listenAddr)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// detecting app termination using signals
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// sigChan channel blocks until its written to
	sig := <-sigChan
	fmt.Println("Recieved terminate, graceful shutdown", sig)

	// graceful shutdown 2 second after termination
	tc, _ := context.WithDeadline(context.Background(), <-time.After(2*time.Second))
	server.Shutdown(tc)

}
