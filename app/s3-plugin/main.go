package main

import (
	"fmt"
	"github.com/nik/Imagitics/platform-s3-plugin/client"
)

func main() {
	bigBuff := make([]byte, 750)
	aws_access_key_id := ""
	aws_secret_access_key := ""
	s3Service, err := client.NewS3Service(aws_access_key_id, aws_secret_access_key, "ap-south-1", "")
	if err == nil {
		s3Location, err := s3Service.Upload("tenant-id-pkg1234", "path", bigBuff)
		if err == nil {
			fmt.Println(s3Location)
		} else {
			fmt.Println(err)
		}
	}
}
