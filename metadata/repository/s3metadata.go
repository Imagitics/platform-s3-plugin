package repository

import "github.com/nik/Imagitics/platform-s3-plugin/pkg/model"

type S3Metadata interface {
	Get(tenantID string) (*model.S3Metadata, error)
}
