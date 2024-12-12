package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileUploader struct {
    client *minio.Client
    bucketName string
}

// NewFileUploader creates a new FileUploader instance
func NewFileUploader() (*FileUploader, error) {
    endpoint := os.Getenv("MINIO_ENDPOINT")
    accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
    secretAccessKey := os.Getenv("MINIO_SECRET_KEY")

    minioClient, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
        Secure: true,
    })
    if err != nil {
        return nil, err
    }

    return &FileUploader{
        client: minioClient,
        bucketName: os.Getenv("MINIO_BUCKET"), // Use environment variable for bucket name
    }, nil
}

// UploadFile uploads a file to MinIO and returns the file URL
func (fu *FileUploader) UploadFile(ctx context.Context, file multipart.File, prefix string) (string, error) {
    // Generate unique object name
    objectName := prefix + "/" + uuid.New().String()

    // Upload the file
    _, err := fu.client.PutObject(ctx, fu.bucketName, objectName, file, -1, minio.PutObjectOptions{
        ContentType: "application/octet-stream", // You can make this more specific
    })
    if err != nil {
        return "", err
    }

    // Construct and return the file URL
    endpoint := os.Getenv("MINIO_ENDPOINT")
    return fmt.Sprintf("%s/%s/%s", endpoint, fu.bucketName, objectName), nil
}

// CreateBucketIfNotExists ensures the bucket exists
func (fu *FileUploader) CreateBucketIfNotExists(ctx context.Context, location string) error {
    exists, err := fu.client.BucketExists(ctx, fu.bucketName)
    if err != nil {
        return err
    }

    if !exists {
        err = fu.client.MakeBucket(ctx, fu.bucketName, minio.MakeBucketOptions{Region: location})
        if err != nil {
            return err
        }
    }

    return nil
}