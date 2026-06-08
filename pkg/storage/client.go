package storage

import "github.com/aws/aws-sdk-go-v2/service/s3"

type S3Client struct {
	Client        *s3.Client
	PresignClient *s3.PresignClient
	Bucket        string
}
