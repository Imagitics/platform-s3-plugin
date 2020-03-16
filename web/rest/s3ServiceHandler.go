package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nik/Imagitics/platform-s3-plugin/client"
	"github.com/nik/Imagitics/platform-s3-plugin/metadata/repository"
	"github.com/nik/Imagitics/platform-s3-plugin/web/rest/model"
	"io/ioutil"
	"net/http"
	"os"
)

const FileSizeLimitError = "multipart: NextPart: http: request body too large"
const NoSuchFileError = "http: no such file"

type S3FileHandler struct {
	repo repository.CassandraS3MetadataRepo
}

func NewS3FileHandler(repoInstance *repository.CassandraS3MetadataRepo) *S3FileHandler {
	s3Handler := &S3FileHandler{
		repo: *repoInstance,
	}

	return s3Handler
}
func S3RegistrationHander(w http.ResponseWriter, r *http.Request) {

}

// S3UploadHandler handles the incoming rest request for post service
// it retrieves bucket, tenant_id and actual entity to be uploaded
// in case of any error it simply returns the error and relevant status code
func (handler *S3FileHandler) S3UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve path parameters
	vars := mux.Vars(r)
	tenantId := vars["tenant_id"]

	//validateTenantId - to be added

	// Retrieve aws credentials for this tenant
	handler.getAWSCredentialsByTenantId(tenantId)

	// Set file limit to configurable size
	r.Body = http.MaxBytesReader(w, r.Body, 2*1024*1024) // 2 Mb
	// Unmarshal the request body into struct and then perform the option of upload
	// Retrieve request and file from the form

	request := r.FormValue("request")
	file, header, err := r.FormFile("entity")

	if err != nil {
		switch err.Error() {
		case FileSizeLimitError:
			http.Error(w, "Exceeded allowed file size limit of 2 Mb", http.StatusBadRequest)
			return
		case NoSuchFileError:
			http.Error(w, "Bad file", http.StatusBadRequest)
			return
		}

		defer file.Close()
		// read all of the contents of our uploaded file into a byte array
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}

		if file == nil {
			fmt.Println("File to be uploaded to s3 can not be empty")
		}

		// Create a temporary file within our temp-images directory that follows
		// a particular naming pattern
		tempFile, err := ioutil.TempFile("", header.Filename)
		tempFile.Write(fileBytes)

		if err != nil {
			fmt.Println(err)
		}

		// This will remove the file once its uploaded to s3 bucket
		defer os.Remove(tempFile.Name())

		s3UploadRequest := model.S3UploadRequest{}
		err = json.Unmarshal([]byte(request), &s3UploadRequest)
		if err == nil {
			w.WriteHeader(402)
		}
		//s3UploadRequest.File = file

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(s3UploadRequest.Bucket)
	}
}

func (handler *S3FileHandler) getAWSCredentialsByTenantId(tenantId string) (*client.S3Credentials, error) {
	if tenantId == "" {
		return nil, errors.New("Tenant-ID can not be empty")
	}

	s3Metadata := handler.repo.Get(tenantId)
	fmt.Println(s3Metadata.AccessKey)
	fmt.Println(s3Metadata.SecretKey)
	fmt.Println(s3Metadata.Region)
	fmt.Println(s3Metadata.TenantID)

	s3Credentials := &client.S3Credentials{
		AWSSecretKey: s3Metadata.SecretKey,
		AWSAccessKey: s3Metadata.AccessKey,
	}

	return s3Credentials, nil
}
