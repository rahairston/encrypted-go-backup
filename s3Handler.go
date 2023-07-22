package main

import (
	"bytes"
	"context"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type BucketHandler struct {
	client *s3.Client
	bucket *string
	prefix string
}

func BuildBucket(bucketName string, keyPrefix string) (*BucketHandler, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return nil, err
	}

	return &BucketHandler{
		client: s3.NewFromConfig(cfg),
		bucket: aws.String(bucketName),
		prefix: keyPrefix,
	}, nil

}

func (bucket BucketHandler) putObject(key string, body []byte) error {

	_, err := bucket.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: bucket.bucket,
		Key:    aws.String(bucket.prefix + "/" + key),
		Body:   bytes.NewReader(body),
	})

	return err
}

func (bucket BucketHandler) getObject(key *string) ([]byte, error) {
	result, err := bucket.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: bucket.bucket,
		Key:    key,
	})

	if err != nil {
		log.Fatal(err)
	}

	return io.ReadAll(result.Body)
}
