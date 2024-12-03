package clients

import "GoBack/types"

type StorageClient interface {
	CheckBucketConnection() bool
	UploadFile(progParams types.ProgParams, fileDiff types.FileDiff) bool
}
