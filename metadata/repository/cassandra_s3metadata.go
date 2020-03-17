package repository

import (
	"errors"
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

// GetMetadataByTenantID retrieves aws metadata for a provided tenant identifier.
// This metadata includes secret key, access key and preferred_region.
func (repo *CassandraS3MetadataRepo) Get(tenantID string) (*model.S3Metadata, error) {
	// Query to retrieve metadata from aws_metadata table
	selectQuery := "select access_key,preferred_region,secret_key from aws_metadata where tenant_id = ?"
	iter := repo.session.Query(selectQuery, tenantID).Iter()
	if iter.NumRows() != 1 {
		//maximum one record is expected as tenant identifier is the unique key
		return nil, errors.New("Bad request")
	}

	// Scan and store relevant attributes into struct
	var s3Metadata model.S3Metadata
	m := map[string]interface{}{}
	iter.MapScan(m)
	s3Metadata = model.S3Metadata{
		AccessKey: m["access_key"].(string),
		Region:    m["preferred_region"].(string),
		SecretKey: m["secret_key"].(string),
		TenantID:  tenantID,
	}

	return &s3Metadata, nil
}
