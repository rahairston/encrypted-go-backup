package aws

import (
	"bytes"
	"context"
	"io"
	"log"
	"strings"

	"github.com/rahairston/EncryptedGoBackup/common"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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

	if !strings.HasPrefix(adjustedKey, "/") && !strings.HasSuffix(bucket.s3Config.Prefix, "/") {
		adjustedKey = bucket.s3Config.Prefix + "/" + adjustedKey
	} else if !strings.HasPrefix(adjustedKey, "/") || !strings.HasSuffix(bucket.s3Config.Prefix, "/") {
		adjustedKey = bucket.s3Config.Prefix + adjustedKey
	} else {
		adjustedKey = bucket.s3Config.Prefix + strings.TrimPrefix(adjustedKey, "/")
	}

	_, err := bucket.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:       &bucket.s3Config.Bucket,
		Key:          aws.String(adjustedKey),
		Body:         bytes.NewReader(body),
		StorageClass: bucket.getTier(adjustedKey),
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

func (bucket BucketHandler) getTier(key string) types.StorageClass {
	for _, element := range bucket.s3Config.Tier.Files { // File match first then folder match
		for _, match := range element.Matches {
			if strings.HasSuffix(key, match) {
				return element.Tier
			}
		}
	}

	for _, element := range bucket.s3Config.Tier.Folders {
		for _, match := range element.Matches {
			folderName := match
			if !strings.HasSuffix(folderName, "/") {
				folderName = folderName + "/" // guarantee folder match and not file match
			}
			if strings.Contains(key, folderName) {
				return element.Tier
			}
		}
	}

	return bucket.s3Config.Tier.Default
}
