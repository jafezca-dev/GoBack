package types

type ProgParams struct {
	BackupType, BackupDate, BasePath, StorageClient, Bucket, Endpoint, AccessKey, SecretKey, Date string
	UploadThread                                                                                  int
	IgnoreFolders, IgnoreFiles                                                                    map[string]bool
}
