package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/phillipmugisa/go_user_api/api"
	"github.com/phillipmugisa/go_user_api/storage"
)

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	listenAddr := flag.String("listenAddr", fmt.Sprintf(":%s", port), "Api Server Port")
	flag.Parse()

	store, err := storage.NewMySqlStorage()
	if err != nil {
		log.Fatal(err)
	}
	err = store.SetUpDB()
	if err != nil {
		log.Fatal(err)
	}

	s := api.NewApiServer(*listenAddr, store)
	log.Fatal(s.Run())
}
