package main

import (
	"github.com/gorilla/mux"
	"github.com/nik/Imagitics/platform-s3-plugin/infra/cassandra"
	"github.com/nik/Imagitics/platform-s3-plugin/metadata/repository"
	"github.com/nik/Imagitics/platform-s3-plugin/web/rest"
	"log"
	"net/http"
	"time"
)

func main() {
	router := mux.NewRouter()

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}


	conn := &cassandra.CassandraConn{
		Hosts:        [] string {"172.18.0.2","172.18.0.2"},
		Port:        "9042",
		User:        "cassandra",
		Password:    "cassandra",
		Consistency: "Quorum",
		Keyspace:    "platform_s3_db",
	}
	repoInstance := repository.NewCassandraS3MetadataRepo(conn)
	handler := rest.NewS3FileHandler(repoInstance)

	router.HandleFunc("/{tenant_id}/s3-store", handler.S3UploadHandler).Methods("POST")

	log.Fatal(srv.ListenAndServe())

}
