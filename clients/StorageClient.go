package clients

import (
	"GoBack/types"
	"bytes"
)

type StorageClient interface {
	CheckBucketConnection() bool
	UploadFile(progParams types.ProgParams, fileDiff types.FileDiff) bool
	UploadCsv(backupInfo bytes.Buffer) bool
}
