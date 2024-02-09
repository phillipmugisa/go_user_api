package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/phillipmugisa/go_user_api/api"
	"github.com/phillipmugisa/go_user_api/storage"
)

func main() {
	port := os.Getenv("PORT")
	listenAddr := flag.String("listenAddr", fmt.Sprintf(":%s", port), "Api Server Port")
	flag.Parse()

	store, err := storage.NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}

	err = store.SetUpDB()
	if err != nil {
		log.Fatal(err)
	}

	s := api.NewApiServer(*listenAddr, store)
	s.Run()
}
