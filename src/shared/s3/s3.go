package s3

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	logger *slog.Logger
	once   sync.Once
)

func InitPackage(_logger *slog.Logger) {
	once.Do(func() {
		logger = _logger
	})
}

// Create MinIO client (reusable)
func NewMinioClient(endpoint, accessKey, secretKey string, useSSL bool) (*minio.Client, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		if logger != nil {
			logger.Error("Failed to create MinIO client",
				slog.String("endpoint", endpoint),
				slog.Bool("useSSL", useSSL),
				slog.Any("error", err),
			)
		}
		return nil, err
	}
	if logger != nil {
		logger.Info("MinIO client created",
			slog.String("endpoint", endpoint),
			slog.Bool("useSSL", useSSL),
		)
	}
	return client, nil
}

// Store object
func StoreObject(ctx context.Context, client *minio.Client, bucketName, objectName, contentType string, data []byte) error {
	if logger != nil {
		logger.Debug("Storing object to S3",
			slog.String("bucket", bucketName),
			slog.String("object", objectName),
			slog.String("contentType", contentType),
			slog.Int("size", len(data)),
		)
	}
	_, err := client.PutObject(ctx, bucketName, objectName,
		bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		if logger != nil {
			logger.Error("Failed to store object",
				slog.String("bucket", bucketName),
				slog.String("object", objectName),
				slog.Any("error", err),
			)
		}
		return err
	}
	if logger != nil {
		logger.Info("Object stored successfully",
			slog.String("bucket", bucketName),
			slog.String("object", objectName),
		)
	}
	return nil
}

// GetObject fetches an object as bytes from MinIO/S3
func GetObject(ctx context.Context, client *minio.Client, bucketName, objectName string) ([]byte, error) {
	if logger != nil {
		logger.Debug("Fetching object from S3",
			slog.String("bucket", bucketName),
			slog.String("object", objectName),
		)
	}
	obj, err := client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		if logger != nil {
			logger.Error("Failed to fetch object (GetObject call failed)",
				slog.String("bucket", bucketName),
				slog.String("object", objectName),
				slog.Any("error", err),
			)
		}
		return nil, err
	}
	defer obj.Close()

	var buf bytes.Buffer
	n, err := io.Copy(&buf, obj)
	if err != nil {
		if logger != nil {
			logger.Error("Failed to read object stream",
				slog.String("bucket", bucketName),
				slog.String("object", objectName),
				slog.Int64("readBytes", n),
				slog.Any("error", err),
			)
		}
		return nil, err
	}
	if logger != nil {
		logger.Info("Object fetched successfully",
			slog.String("bucket", bucketName),
			slog.String("object", objectName),
			slog.Int("size", buf.Len()),
		)
	}
	return buf.Bytes(), nil
}
