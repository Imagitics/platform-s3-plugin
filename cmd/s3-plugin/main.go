package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nik/Imagitics/platform-s3-plugin/infra/cassandra"
	"github.com/nik/Imagitics/platform-s3-plugin/metadata/repository"
	"github.com/nik/Imagitics/platform-s3-plugin/web/api/v1"
	"log"
	"net/http"
	"time"
)

func main() {
	//create router for routing the requests
	router := mux.NewRouter()

	//create cassandra connection instance
	conn := &cassandra.CassandraConn{
		Hosts:       []string{"172.18.0.2"},
		Port:        "9042",
		User:        "cassandra",
		Password:    "cassandra",
		Consistency: "Quorum",
		Keyspace:    "platform_s3_db",
	}

	//create repoinstance
	repoInstance := repository.NewCassandraS3MetadataRepo(conn)
	fmt.Println("Connected to Cassandra")

	//instantiate api object and routes
	apiInstnace := v1.NewApi(router, repoInstance)
	apiInstnace.InitializeRoutes()

	//create a http server
	srv := &http.Server{
		Addr: "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler:      router,
	}
	fmt.Println("Initializing http server")
	log.Fatal(srv.ListenAndServe())
}
