package model

type S3Metadata struct {
	tenantID  string `json:"tenantID"`
	accessKey string `json:"accessKey"`
	secretKey string `json:"secretKe"`
	region    string `json:"region"`
}
