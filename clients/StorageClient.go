package clients

import (
	"GoBack/types"
	"bytes"
)

type StorageClient interface {
	CheckBucketConnection() bool
	UploadFile(fileDiff types.FileDiff) bool
	UploadCsv(backupInfo bytes.Buffer) bool
	GetLastChanges() (map[string]types.FileChanges, error)
	MultiThreadUpload([]types.FileDiff) bool
}
