package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nik/Imagitics/platform-s3-plugin/infra/cassandra"
	"github.com/nik/Imagitics/platform-s3-plugin/metadata/repository"
	"github.com/nik/Imagitics/platform-s3-plugin/utility"
	"github.com/nik/Imagitics/platform-s3-plugin/web/api/v1"
	"log"
	"net/http"
	"time"
)

func main() {
	//load configuration
	config, err := utility.LoadConfiguration("/etc/config/config.json")
	if err != nil {
		//halt bootstrapping
		fmt.Println("Error in loading configuration - ", err)
	}

	fmt.Println(config.Cassandra.Host)
	fmt.Println(config.Cassandra.Port)

	//instantiate cassandra connection instance
	conn := &cassandra.CassandraConn{
		Hosts:       []string{config.Cassandra.Host},
		Port:        config.Cassandra.Port,
		User:        config.Cassandra.User,
		Password:    config.Cassandra.Password,
		Consistency: config.Cassandra.Consistency,
		Keyspace:    config.Cassandra.Keyspace,
		CaPath:      config.Cassandra.SSLCertPath,
	}

	//create repoinstance
	repoInstance := repository.NewCassandraS3MetadataRepo(conn)
	fmt.Println("Connected to Cassandra")

	//instantiate api object and routes
	router := mux.NewRouter()
	apiInstnace := v1.NewApi(router, repoInstance)
	apiInstnace.InitializeRoutes()

	//create a http server
	srv := &http.Server{
		Addr: ":8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		Handler:      router,
	}
	fmt.Println("Initializing http server")
	log.Fatal(srv.ListenAndServe())
}
