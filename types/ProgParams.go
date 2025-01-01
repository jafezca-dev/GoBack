package types

type ProgParams struct {
	BackupType, BackupDate, BasePath, StorageClient, Bucket, Endpoint, AccessKey, SecretKey string
	IgnoreFolders, IgnoreFiles                                                              map[string]bool
}
