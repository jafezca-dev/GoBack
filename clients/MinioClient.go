package clients

import (
	"GoBack/types"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client     minio.Client
	ProgParams types.ProgParams
}

func NewMinioClient(progParams types.ProgParams) *MinioClient {
	minioClient, err := minio.New(progParams.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(progParams.AccessKey, progParams.SecretKey, ""),
		Secure: true,
	})

	if err != nil {
		panic("Error creating minio client")
	}

	return &MinioClient{
		ProgParams: progParams,
		client:     *minioClient,
	}
}

func (mc *MinioClient) CheckBucketConnection() bool {
	_, err := mc.client.GetBucketPolicy(context.Background(), mc.ProgParams.Bucket)
	if err != nil {
		panic("Error checking bucket connection")
	}

	return true
}

func (mc *MinioClient) UploadFile(progParams types.ProgParams, fileDiff types.FileDiff) bool {
	_, virtualFilePath := fileDiff.DirPaths(progParams.BasePath)
	bucketPath := progParams.BackupDate + virtualFilePath
	bucketPath = strings.ReplaceAll(bucketPath, "\\", "/")
	fmt.Println(bucketPath)
	contentType := "text/plain"

	file, err := os.Open(fileDiff.FullPath())
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	_, err = mc.client.PutObject(context.Background(), mc.ProgParams.Bucket, bucketPath, file, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})

	// fmt.Println(info)
	return true
}

func (mc *MinioClient) UploadCsv(backupInfo bytes.Buffer) bool {
	_, _ = mc.client.PutObject(context.Background(), mc.ProgParams.Bucket, "/backup_data/"+mc.ProgParams.BackupDate+".csv",
		bytes.NewReader(backupInfo.Bytes()), int64(backupInfo.Len()), minio.PutObjectOptions{
			ContentType: "text/csv",
		})

	return true
}
