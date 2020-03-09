package client

import(
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
"github.com/aws/aws-sdk-go/aws/credentials"
"github.com/aws/aws-sdk-go/service/s3"
"github.com/aws/aws-sdk-go/aws/session"
"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

)

type S3FileUploadRequest struct {
	Bucket string
	DirectoryPath string
	file []byte
}

func Upload(bucketName string, directoryPath string, file []byte ) {
	aws_access_key_id := ""
	aws_secret_access_key := ""
	token := ""
	creds := credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, token)

	_, err := creds.Get()
	if err != nil {
		// handle error
	}
	cfg := aws.NewConfig().WithRegion("us-west-1").WithCredentials(creds)

	session, err:= session.NewSession(cfg)

		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}


	s3Client := s3.New(session)
	bucketInstance:= &s3.CreateBucketInput{Bucket:aws.String(bucketName)}
	result, err:= s3Client.CreateBucket(bucketInstance)
	println(result.Location)
	if(err!=nil) {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists
				goto next
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}

	next:
		// Create an uploader with the session and default options
		uploader := s3manager.NewUploader(session)

		// Upload the file to S3.
		fileUploadReuslt, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(directoryPath),
			Body:   file,
		})

	fmt.Println(fileUploadReuslt.Location)
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)


}
