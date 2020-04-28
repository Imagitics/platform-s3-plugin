package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nik/Imagitics/platform-s3-plugin/client"
	model2 "github.com/nik/Imagitics/platform-s3-plugin/metadata/model"
	"github.com/nik/Imagitics/platform-s3-plugin/metadata/repository"
	"github.com/nik/Imagitics/platform-s3-plugin/web/rest/model"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

const FileSizeLimitError = "multipart: NextPart: http: request body too large"
const NoSuchFileError = "http: no such file"

type Api struct {
	router *mux.Router
	repo repository.CassandraS3MetadataRepo
}

func NewApi(router *mux.Router,repoInstance *repository.CassandraS3MetadataRepo) *Api {
	s3Handler := &Api{
		repo: *repoInstance,
		router:router,
	}

	return s3Handler
}
func S3RegistrationHander(w http.ResponseWriter, r *http.Request) {

}

// S3UploadHandler handles the incoming rest request for post service
// it retrieves bucket, tenant_id and actual entity to be uploaded
// in case of any error it simply returns the error and relevant status code
func (api *Api) upload (w http.ResponseWriter, r *http.Request) {
	// Set file limit to configurable size
	r.Body = http.MaxBytesReader(w, r.Body, 2*1024*1024) // 2 Mb

	// Unmarshal the request body into struct and then perform the option of upload
	// Retrieve request and file from the form
	s3UploadRequest, err := validateAndRetriveUploadRequest(r.FormValue("request"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	//retrieve aws metadata stored agains tenant-id
	s3Credentials, s3Metadata, err := api.getAWSCredentialsAndMetadataByTenantId(s3UploadRequest.TenantId)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Validate the uploaded file and retrieve the file handler for uploading the physical file to s3
	tempFile, err := validateAndGetFileHandler(r.FormFile("entity"))
	if err != nil {
		// File validation failed
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// All validations are passed. Now a file can be uploaded to s3.
	// This will remove the file once its uploaded to s3 bucket
	defer os.Remove(tempFile.Name())

	// Populate region data.It has to be used from the upload request or from the metadata.
	region := s3Metadata.Region
	if s3UploadRequest.Region == "" {
		region = s3UploadRequest.Region
	}
	if region == "" {
		//region is empty
		http.Error(w, "Bad request as region is missing", http.StatusBadRequest)
	}

	//Create a new S3Service instance and use this handler to perform upload operation to s3 bucket
	s3Service, err := client.NewS3Service(s3Credentials.AWSAccessKey, s3Credentials.AWSSecretKey, region, "")
	if err == nil {
		fileBytes, _ := ioutil.ReadAll(tempFile)
		s3Location, err := s3Service.Upload(s3UploadRequest.Bucket, s3UploadRequest.TenantId, fileBytes)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("File can not be uploaded")))
			return
		}
		w.WriteHeader(201)
		w.Write([]byte(fmt.Sprintf("File successfully updated at location %s", s3Location)))
	}
}

// getAWSCredentialsByTenantId retrieves aws credentials for the provided tenant identifier.
func (api *Api) getAWSCredentialsAndMetadataByTenantId(tenantId string) (*client.S3Credentials, *model2.S3Metadata, error) {
	if tenantId == "" {
		// tenantId can not be empty.Its better that the validation is done at a higher level
		return nil, nil, errors.New("Tenant-ID can not be empty")
	}

	// Retrieve aws metadata
	s3Metadata, err := api.repo.Get(tenantId)
	if err != nil {
		return nil, nil, err
	}
	// Populate data structure as per the attributes
	s3Credentials := &client.S3Credentials{
		AWSSecretKey: s3Metadata.SecretKey,
		AWSAccessKey: s3Metadata.AccessKey,
	}

	return s3Credentials, s3Metadata, nil
}

// validateAndGetFileHandler performs validation on uploaded file.
// File can not be nil and there is a limit on the file size.
// If all file validations are passed, it creates temporary file at /tmp location
// This physical file (stored on disk) is to be used later to upload
func validateAndGetFileHandler(file multipart.File, header *multipart.FileHeader, err error) (*os.File, error) {
	if file == nil {
		// file object must not be nil
		return nil, errors.New("File to be uploaded to s3 can not be empty")
	}

	if err != nil {
		// check for error and corresponding cause
		switch err.Error() {
		case FileSizeLimitError:
			return nil, errors.New("Exceeded allowed file size limit of 2 Mb")
		case NoSuchFileError:
			return nil, errors.New("Bad file")
		}
	}

	// Read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	// create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("", header.Filename)
	if err != nil {
		return nil, err
	}
	tempFile.Write(fileBytes)

	return tempFile, nil
}

// validateAndRetriveUploadRequest validates the uploaded file upload request
// It retrieves tenant identifier and bucket attributes from the request
func validateAndRetriveUploadRequest(body string) (*model.S3UploadRequest, error) {
	// initialize struct instance and try to retrieve request object from the input form paramter
	s3UploadRequest := &model.S3UploadRequest{}
	err := json.Unmarshal([]byte(body), &s3UploadRequest)
	if err != nil {
		//error in unmarshaling because of format mismatch
		return nil, err
	}

	if s3UploadRequest.Bucket == "" {
		//validate empty bucket case
		return nil, errors.New("Bucket can not be empty")
	} else if s3UploadRequest.TenantId == "" {
		//validate tenant identifier case
		return nil, errors.New("Invalid request")
	}

	return s3UploadRequest, nil
}

func (a *Api) InitializeRoutes() {
	//a.router.HandleFunc("/s3/images/{id:[0-9]+}", a.getProducts).Methods("GET")
	a.router.HandleFunc("/s3/files", a.upload).Methods("POST")
	//a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	//a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	//a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
}