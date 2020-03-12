package main

import (
	"fmt"
	"github.com/nik/Imagitics/platform-s3-plugin/client"
	"github.com/nik/Imagitics/platform-s3-plugin/utility"
	"path/filepath"
)

func main() {
	bigBuff := make([]byte, 750)
	configPath := filepath.FromSlash("/app/platform-s3-plugin/config/config.json")

	config, error := utility.LoadConfiguration(configPath)
	if error != nil {
		panic(error)
	}

	s3Service, err := client.NewS3Service(config.S3.AwsAccessKey, config.S3.AwsSecretKey, config.S3.Region, "")
	if err == nil {
		s3Location, err := s3Service.Upload("tenant-id-pkg1234", "path", bigBuff)
		if err == nil {
			fmt.Println(s3Location)
		} else {
			fmt.Println(err)
		}
	}
}
