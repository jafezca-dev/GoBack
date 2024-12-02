package clients

type StorageClient interface {
	CheckBucketConnection() bool
}
