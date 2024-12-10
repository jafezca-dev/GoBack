package clients

import (
	"GoBack/types"
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

func (mc *MinioClient) GetLastChanges() (map[string]types.FileChanges, error) {
	prefix := "backup_data/"

	backupDataFiles := mc.client.ListObjects(context.Background(), mc.ProgParams.Bucket, minio.ListObjectsOptions{
		Prefix:    "backup_data/",
		Recursive: false,
	})

	responseLen := 0
	var latestFile minio.ObjectInfo

	for file := range backupDataFiles {
		responseLen += 1
		if file.Err != nil {
			log.Fatalln(file.Err)
		}

		if file.Key == prefix || file.Key[len(file.Key)-1] == '/' {
			continue
		}

		if file.LastModified.After(latestFile.LastModified) {

			latestFile = file
		}
	}

	if responseLen == 0 {
		return nil, fmt.Errorf("no changes found")
	}

	file, err := mc.client.GetObject(context.Background(), mc.ProgParams.Bucket, latestFile.Key, minio.GetObjectOptions{})
	defer file.Close()

	if err != nil {
		return nil, fmt.Errorf("Error getting latest changes for backup data")
	}

	scanner := bufio.NewScanner(file)

	var changesData []string

	for scanner.Scan() {
		changesData = append(changesData, scanner.Text())
	}

	//Create new struct
	oldData := make(map[string]types.FileChanges)

	dateFormat := "2006-01-02 15:04:05"

	for _, change := range changesData {
		//fmt.Println(change)
		splitedString := strings.Split(change, ";")

		formatedTime, _ := time.Parse(dateFormat, splitedString[1])
		oldData[splitedString[0]] = types.FileChanges{ModTime: formatedTime, BackupTag: splitedString[2]}
	}

	return oldData, nil
}
