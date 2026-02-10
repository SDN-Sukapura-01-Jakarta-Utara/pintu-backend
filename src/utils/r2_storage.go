package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// R2Storage handles all R2 operations
type R2Storage struct {
	client     *s3.Client
	bucketName string
	publicURL  string
}

// NewR2Storage initializes R2 storage
func NewR2Storage() *R2Storage {
	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("R2_BUCKET_NAME")
	endpoint := os.Getenv("R2_ENDPOINT")
	publicDomain := os.Getenv("R2_PUBLIC_DOMAIN")

	// Create credentials
	creds := credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")

	// Create S3 client with R2 endpoint
	client := s3.NewFromConfig(aws.Config{
		Region:      "auto",
		Credentials: creds,
		BaseEndpoint: aws.String(endpoint),
	})

	return &R2Storage{
		client:     client,
		bucketName: bucketName,
		publicURL:  publicDomain,
	}
}

// UploadFile uploads file to R2 storage
func (r *R2Storage) UploadFile(file *multipart.FileHeader, directory string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Open file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Read file content
	fileContent, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Generate filename with timestamp (inside directory)
	timestamp := time.Now().Unix()
	sanitizedFilename := sanitizeFilename(file.Filename)
	filename := fmt.Sprintf("%s/%d-%s", directory, timestamp, sanitizedFilename)

	// Upload to R2
	putObjectInput := &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(fileContent),
		ContentType: aws.String(file.Header.Get("Content-Type")),
	}

	_, err = r.client.PutObject(ctx, putObjectInput)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to R2: %w", err)
	}

	// Return the file path (object key in R2)
	return filename, nil
}

// DeleteFile deletes file from R2 storage
func (r *R2Storage) DeleteFile(fileKey string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	deleteObjectInput := &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(fileKey),
	}

	_, err := r.client.DeleteObject(ctx, deleteObjectInput)
	if err != nil {
		return fmt.Errorf("failed to delete file from R2: %w", err)
	}

	return nil
}

// GetPublicURL returns the public URL for a file
func (r *R2Storage) GetPublicURL(fileKey string) string {
	// Format: https://pintu-storage.sdnsukapura01.sch.id/<fileKey>
	return fmt.Sprintf("https://%s/%s", r.publicURL, fileKey)
}

// GetPresignedURL returns a presigned URL for a file (optional, for private files)
func (r *R2Storage) GetPresignedURL(fileKey string, expirationMinutes int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create presigner from client
	presigner := s3.NewPresignClient(r.client)

	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(fileKey),
	}

	presignedRequest, err := presigner.PresignGetObject(ctx, getObjectInput,
		func(opts *s3.PresignOptions) {
			opts.Expires = time.Duration(expirationMinutes) * time.Minute
		})

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedRequest.URL, nil
}

// sanitizeFilename removes spaces and special characters from filename
func sanitizeFilename(filename string) string {
	// Replace spaces with dash
	filename = regexp.MustCompile(`\s+`).ReplaceAllString(filename, "-")
	
	// Remove special characters, keep only alphanumeric, dash, dot, underscore
	filename = regexp.MustCompile(`[^a-zA-Z0-9.-_]`).ReplaceAllString(filename, "")
	
	return filename
}
