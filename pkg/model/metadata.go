package model

type S3Metadata struct {
	TenantID  string `json:"tenantID"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKe"`
	Region    string `json:"region"`
}
