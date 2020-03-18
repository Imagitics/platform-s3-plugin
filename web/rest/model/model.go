package model

type S3UploadRequest struct {
	Bucket   string `json:"bucket"`
	TenantId string `json:"tenant_id"`
	File     []byte `json:"entity`
	Region   string `json:"region`
}
