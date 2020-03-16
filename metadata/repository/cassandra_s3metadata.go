package repository

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/nik/Imagitics/platform-s3-plugin/infra/cassandra"
	"github.com/nik/Imagitics/platform-s3-plugin/metadata/model"
)

type CassandraS3MetadataRepo struct {
	session *gocql.Session
}

func NewCassandraS3MetadataRepo(conn *cassandra.CassandraConn) *CassandraS3MetadataRepo {
	conn.Keyspace = "platform_s3_db"
	conn.Consistency = "QUORUM"
	session := conn.InitSession()
	repo := &CassandraS3MetadataRepo{
		session: session,
	}

	return repo
}

func (repo *CassandraS3MetadataRepo) Get(tenantID string) {
	selectQuery := "select access_key,preferred_region,secret_key from aws_metadata where tenant_id = ?"
	iter := repo.session.Query(selectQuery, tenantID).Iter()
	var s3Metadata model.S3Metadata
	iter.Scan(&s3Metadata)

	fmt.Println(s3Metadata.Region)
}
