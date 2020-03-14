package main

import (
	"github.com/nik/Imagitics/platform-s3-plugin/web/rest"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/s3-store", rest.S3UploadHandler).Methods("POST")

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
