package storage

import (
	"context"
	"io"
	"net/http"
	"path"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

type S3Storage struct {
	client *s3.Client
	bucket string
	url    string
}

func (s *S3Storage) Upload(r io.Reader, fpath string) string {
	out, err := s.client.PutObject(
		context.TODO(),
		&s3.PutObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(path),
			Body:   r,
		},
	)

	return path.Join(url, path)
}

func (s *S3Storage) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
