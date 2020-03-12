package config

type ConfigModel struct {
	S3 struct {
		AwsAccessKey string `json:"aws_access_key"`
		AwsSecretKey string `json:"aws_secret_key"`
		Region       string `json:"region"`
		PathPrefix   string `json:"path_prefix"`
	} `json:"s3"`
}