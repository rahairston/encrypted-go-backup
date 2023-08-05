package aws

import (
	"bytes"
	"context"
	"io"
	"log"
	"strings"

	"backup/common"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BucketHandler struct {
	client   *s3.Client
	s3Config *common.S3Object
}

func BuildBucket(backupConfig *common.BackupConfig) (*BucketHandler, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(backupConfig.Profile))

	if err != nil {
		return nil, err
	}

	return &BucketHandler{
		client:   s3.NewFromConfig(cfg),
		s3Config: &backupConfig.S3Config,
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
		Bucket:       &bucket.s3Config.Bucket,
		Key:          aws.String(bucket.s3Config.Prefix + adjustedKey),
		Body:         bytes.NewReader(body),
		StorageClass: bucket.s3Config.Tier,
	})

	return err
}

func (bucket BucketHandler) GetObject(key string) ([]byte, error) {
	result, err := bucket.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucket.s3Config.Bucket,
		Key:    &key,
	})

	if err != nil {
		log.Fatal(err)
	}

	return io.ReadAll(result.Body)
}
