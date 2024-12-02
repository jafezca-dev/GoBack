package clients

import (
	"GoBack/progparams"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client     minio.Client
	ProgParams progparams.ProgParams
}

func NewMinioClient(progParams progparams.ProgParams) *MinioClient {
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
