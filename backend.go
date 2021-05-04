package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Backend is a generic interface for a backup storage backend
type Backend interface {
	ListFiles(prefix string) ([]*File, error)
	CopyFile(source string, destination string) error
	DeleteFile(path string) error
}

// LocalBackend implements a local fs backend
type LocalBackend struct {
	root string
}

// NewLocalBackend returns a LocalBackend at the specified root dir
func NewLocalBackend(root string) (*LocalBackend, error) {
	return &LocalBackend{root: root}, nil
}

var _ Backend = (*LocalBackend)(nil)

// ListFiles lists all files in the backend at the given path in reverse modified order
func (l *LocalBackend) ListFiles(path string) ([]*File, error) {
	dir, err := os.Open(filepath.Join(l.root, path))
	if err != nil {
		return nil, err
	}

	files, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	ret := make([]*File, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		ret = append(ret, &File{
			Name:    filepath.Join(path, file.Name()),
			Created: file.ModTime(),
		})
	}

	return ret, nil
}

// CopyFile a
func (l *LocalBackend) CopyFile(source string, destination string) error {
	src, err := os.Open(filepath.Join(l.root, source))
	if err != nil {
		return err
	}

	// make destination directory if it doesn't exist
	dir := filepath.Dir(filepath.Join(l.root, destination))
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	dst, err := os.Create(filepath.Join(l.root, destination))
	if err != nil {
		return err
	}

	_, err = io.Copy(src, dst)
	return err
}

// DeleteFile a
func (l *LocalBackend) DeleteFile(file string) error {
	return os.Remove(filepath.Join(l.root, file))
}

// S3Backend a
type S3Backend struct {
	bucket string
	client *s3.Client
	ctx    context.Context
}

// NewS3Backend a
func NewS3Backend(ctx context.Context, bucket string) (*S3Backend, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	return &S3Backend{bucket: bucket, client: client, ctx: ctx}, nil
}

var _ Backend = (*S3Backend)(nil)

// ListFiles a
func (s *S3Backend) ListFiles(prefix string) ([]*File, error) {
	if prefix != "" {
		prefix = prefix + "/"
	}
	ret := make([]*File, 0)
	nextToken := (*string)(nil)

	for {
		res, err := s.client.ListObjectsV2(s.ctx, &s3.ListObjectsV2Input{
			Bucket:            aws.String(s.bucket),
			Prefix:            aws.String(prefix),
			Delimiter:         aws.String("/"),
			ContinuationToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, s3file := range res.Contents {
			ret = append(ret, &File{
				Name:    *s3file.Key,
				Created: *s3file.LastModified,
			})
		}

		nextToken = res.NextContinuationToken
		if nextToken == nil {
			break
		}
	}

	return ret, nil
}

// CopyFile a
func (s *S3Backend) CopyFile(source string, destination string) error {
	_, err := s.client.CopyObject(s.ctx, &s3.CopyObjectInput{
		CopySource: aws.String(fmt.Sprintf("%s/%s", s.bucket, source)),
		Bucket:     aws.String(s.bucket),
		Key:        aws.String(destination),
	})

	return err
}

// DeleteFile a
func (s *S3Backend) DeleteFile(path string) error {
	_, err := s.client.DeleteObject(s.ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})

	return err
}
