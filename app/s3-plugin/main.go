package main

import (
	"github.com/nik/Imagitics/platform-s3-plugin/client"
	"github.com/nik/Imagitics/platform-s3-plugin/store"
)

func main() {
	bigBuff := make([]byte, 750000000)
	client.Upload("tenant-id-pkg","path",bigBuff)
}
