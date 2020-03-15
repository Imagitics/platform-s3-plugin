package main

import "github.com/nik/Imagitics/platform-s3-plugin/infra/cassandra"

func main() {
	conn:= cassandra.CassandraConn{
		Host:"cassandra.us-east-1.amazonaws.com",
		Port: "9142",
		User:"CassandraMS-at-948987282987",
		Password:"HNA4B1211FjKoC7Xb2eetbfC8gyFvEw0WXEM180QT9c=",
		Consistency: "Quorum",
		Keyspace:"platform_s3_db",
	}


	cassandra:= conn.InitSession()
	cassandra.Close()
}
