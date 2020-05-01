package model

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Credentials struct {
	AWSAccessKey string
	AWSSecretKey string
	Token        string
}

type S3Client struct {
	Region string
	Config *aws.Config
}

// Initializes S3 credentials with aws access key and aws secret key
// It returns pointer to the object that wraps the aws credentials
func NewS3Credential(accessKey string, secretKey string, token string) (*S3Credentials, error) {
	if accessKey == "" {
		// AWS Access Key can not be empty
		return nil, errors.New("AWS access key can not be empty")
	}

	if secretKey == "" {
		// AWS Secret Key can not be empty
		return nil, errors.New("AWS secret key can not be empty")
	}

	// Wrap credentials within S3Credentials object
	s3Credentials := S3Credentials{AWSAccessKey: accessKey,
		AWSSecretKey: secretKey,
		Token:        token,
	}

	return &s3Credentials, nil
}

// Initializes s3 client with provided region and credentials
// It returns the pointer to the object that wraps S3Client
func (cred *S3Credentials) NewS3Client(region string) (*s3.S3, *session.Session, error) {
	if region == "" {
		// Region can not be empty
		return nil, nil, errors.New("Region can not be empty")
	}

	// Initialize the S3Credential object with parmaeters
	creds := credentials.NewStaticCredentials(cred.AWSAccessKey, cred.AWSSecretKey, cred.Token)
	// Retrieve credentials value
	_, err := creds.Get()

	if err != nil {
		return nil, nil, err
	}

	// Create configure object to be used to create a new session
	cfg := aws.NewConfig().WithRegion(region).WithCredentials(creds)

	// Create a new Session
	session, err := session.NewSession(cfg)
	if err != nil {
		return nil, nil, err
	}

	// Create a new S3Client with the session created in the previous step
	s3Client := s3.New(session)

	return s3Client, session, nil
}
