package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
)

// Performs upload to s3 bucket using configuration specified in the config file
func s3UploadUsingConfig() {
	creds := credentials.NewStaticCredentials("AKIAIHJZ6JMA6LJ4KQBQ", "pv2/HpVjBFEi4fUGMb6oXuaA6ULh4sH1QgUwKRAZ", "")
	// Retrieve credentials value
	_, err := creds.Get()

	if err != nil {
		os.Exit(0)
	}

	// Create configure object to be used to create a new session
	cfg := aws.NewConfig().WithRegion("ap-south-1").WithCredentials(creds)

	// Create a new Session
	ses, err := session.NewSession(cfg)
	svc := s3.New(ses)

	input := &s3.CreateBucketInput{
		Bucket: aws.String("uwrwer"),
	}

	res, err := svc.CreateBucket(input)
	println(res.Location)
	println(err.Error())

}

func main() {
	s3UploadUsingConfig()
}
