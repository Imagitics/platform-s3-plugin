package service

import (
	"bytes"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/nik/Imagitics/platform-s3-plugin/metadata/repository"
	"github.com/nik/Imagitics/platform-s3-plugin/pkg/model"
)

type S3Service struct {
	s3Client *s3.S3
	session  *session.Session
	repo     repository.S3Metadata
}

// Initializes s3Client with required parameters such as accesskey, secretkey and region
// Token is used for aws session token
// It returns pointer to the object that wraps s3Client
func NewS3Service(tenantID string, region string, repo repository.S3Metadata, token string) (*S3Service, error) {

	//retrieve aws metadata stored agains tenant-id
	s3Credentials, s3Metadata, err := getS3AttributesMetadataByTenantId(tenantID, repo)
	if err != nil {
		return nil, err
	}

	//instantiate aws credentials model
	s3Cred, err := model.NewS3Credential(s3Credentials.AWSAccessKey, s3Credentials.AWSSecretKey, token)
	if err != nil {
		return nil, err
	}

	if region == "" {
		region = s3Metadata.Region
	}

	// Provide region to initialize s3Client based on the supplied s3Credentials
	s3Client, session, err := s3Cred.NewS3Client(s3Metadata.Region)
	if err != nil {
		return nil, err
	}

	// Wrap s3Client withing s3Service object
	s3Service := &S3Service{
		s3Client: s3Client,
		session:  session,
	}

	return s3Service, nil
}

// Uploads file to s3 to provided bucket at input directory path
func (s3Service *S3Service) Upload(bucketName string, directoryPath string, file []byte) (string, error) {
	// Create bucket instance
	bucketInstance := &s3.CreateBucketInput{Bucket: aws.String(bucketName)}
	_, err := s3Service.s3Client.CreateBucket(bucketInstance)
	if err != nil {
		// Check whether bucket already exists.
		// Its possible as we we support atleast once semantics
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				goto next
			case s3.ErrCodeBucketAlreadyExists:
				goto next
			default:
				return "", aerr
			}
		}
	}

next:
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(s3Service.session)
	// Convert file into
	body := bytes.NewReader(file)
	// Upload the file to S3.
	fileUploadReuslt, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(directoryPath),
		Body:   body,
	})

	if err != nil {
		return "", err
	}

	return fileUploadReuslt.Location, nil
}

// getAWSCredentialsByTenantId retrieves aws credentials for the provided tenant identifier.
func getS3AttributesMetadataByTenantId(tenantId string, repo repository.S3Metadata) (*model.S3Credentials, *model.S3Metadata, error) {
	if tenantId == "" {
		// tenantId can not be empty.Its better that the validation is done at a higher level
		return nil, nil, errors.New("Tenant-ID can not be empty")
	}

	// Retrieve aws metadata
	s3Metadata, err := repo.Get(tenantId)
	if err != nil {
		return nil, nil, err
	}
	// Populate data structure as per the attributes
	s3Credentials := &model.S3Credentials{
		AWSSecretKey: s3Metadata.SecretKey,
		AWSAccessKey: s3Metadata.AccessKey,
	}

	return s3Credentials, s3Metadata, nil
}
