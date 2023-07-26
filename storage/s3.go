package storage

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client *s3.Client
	bucket string
	prefix string
	url    *url.URL
}

func NewS3Storage(bucket, prefix, urlstr string) (*S3Storage, error) {
	store := S3Storage{
		bucket: bucket,
		prefix: prefix,
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	store.client = s3.NewFromConfig(cfg)
	if urlstr == "" {
		store.url = &url.URL{
			Scheme: "https",
			Host:   bucket + ".s3." + cfg.Region + ".amazonaws.com",
		}
	} else {
		store.url, err = url.Parse(urlstr)
		if err != nil {
			return nil, err
		}
	}

	return &store, nil
}

func (s *S3Storage) Upload(r io.Reader, fpath string) (string, error) {
	_, err := s.client.PutObject(
		context.TODO(),
		&s3.PutObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(path.Join(s.prefix, fpath)),
			Body:   r,
		},
	)
	if err != nil {
		return "", err
	}

	url := s.url.JoinPath(s.prefix, fpath)
	return url.String(), err
}

func (s *S3Storage) Delete(fpath string) error {
	_, err := s.client.DeleteObject(
		context.TODO(),
		&s3.DeleteObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(path.Join(s.prefix, fpath)),
		},
	)

	return err
}

func (s *S3Storage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// NOP
}

func (s S3Storage) Backend() string {
	return "s3"
}
