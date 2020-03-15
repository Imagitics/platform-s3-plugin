package metadata

type MetadataRepository interface {
	Get(tenantId string) error
	Create(tenantID string, accessKey string, secretKey string, preferredRegion string) error
	Update(tenantID string, accessKey string, secretKey string, preferredRegion string) error
	Delete(tenantID string)
}
