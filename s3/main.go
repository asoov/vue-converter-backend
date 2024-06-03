package s3Utils

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Service struct{}

type S3ServiceInterface interface {
	CreateBucket(bucketName string) error
	UploadToBucket(bucketName string, fileName string, fileContent []byte) error
}

func (s *S3Service) CreateBucket(bucketName string) error {
	svc := s3.New(session.Must(session.NewSession()))

	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket:                    &bucketName,
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{LocationConstraint: aws.String("eu-central-1")},
	})

	return err
}

func (s *S3Service) UploadToBucket(bucketName string, fileName string, fileContent []byte) error {
	svc := s3.New(session.Must(session.NewSession()))

	fileContentReader := bytes.NewReader(fileContent)

	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    aws.String(fileName),
		Body:   fileContentReader,
	})

	return err
}
