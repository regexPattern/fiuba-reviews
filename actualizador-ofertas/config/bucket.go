package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	BucketMaxRequests int    = 5
	bucketEnvVarName  string = "AWS_S3_BUCKET"
)

type S3Config struct {
	Client     *s3.Client
	BucketName *string
}

func newS3Client() (*S3Config, error) {
	bn, ok := os.LookupEnv(bucketEnvVarName)
	if !ok {
		return nil, fmt.Errorf("variable de entorno `%v` no configurada", bucketEnvVarName)
	}
	cfg, err := loadAwsDefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("error cargando la configuración por defecto de AWS: %w", err)
	}
	cl := s3.NewFromConfig(cfg)
	if err := checkBucketConn(cl, bn); err != nil {
		return nil, fmt.Errorf("error estableciendo conexión con el bucket: %w", err)
	}
	c := &S3Config{
		Client:     cl,
		BucketName: &bn,
	}
	return c, nil
}

func loadAwsDefaultConfig() (aws.Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	return config.LoadDefaultConfig(ctx)
}

func checkBucketConn(c *s3.Client, bucket string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	_, err := c.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: &bucket,
	})
	return err
}
