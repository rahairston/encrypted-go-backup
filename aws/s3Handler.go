package aws

import (
	"bytes"
	"context"
	"io"
	"log"
	"strings"

	"backup/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BucketHandler struct {
	client *s3.Client
	bucket *string
	prefix string
}

func BuildBucket(backupConfig *types.BackupConfig) (*BucketHandler, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(backupConfig.Profile))

	if err != nil {
		return nil, err
	}

	return &BucketHandler{
		client: s3.NewFromConfig(cfg),
		bucket: aws.String(backupConfig.Bucket),
		prefix: backupConfig.Prefix,
	}, nil

}

func (bucket BucketHandler) PutObject(key string, body []byte) error {
	var adjustedKey = key
	if strings.Contains(key, "\\") { // AWS uses Linux-like pathing
		adjustedKey = strings.Replace(key, "\\", "/", -1)
	}

	if !strings.HasSuffix(adjustedKey, "/") {
		adjustedKey = "/" + adjustedKey
	}

	_, err := bucket.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: bucket.bucket,
		Key:    aws.String(bucket.prefix + adjustedKey),
		Body:   bytes.NewReader(body),
	})

	return err
}

func (bucket BucketHandler) GetObject(key string) ([]byte, error) {
	result, err := bucket.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: bucket.bucket,
		Key:    aws.String(key),
	})

	if err != nil {
		log.Fatal(err)
	}

	return io.ReadAll(result.Body)
}
