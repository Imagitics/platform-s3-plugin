package main

import "github.com/nik/Imagitics/platform-s3-plugin/infra/cassandra"

func main() {
	conn := cassandra.CassandraConn{
		Hosts:        [] string {"172.18.0.2"},
		Port:        "9042",
		User:        "cassandra",
		Password:    "cassandra",
		Consistency: "Quorum",
		Keyspace:    "platform_s3_db",
	}

	cassandra := conn.InitSession()
	cassandra.Close()
}
