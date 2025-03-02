package aws

import (
	"context"
	"fmt"
	internalConf "quiz-app-be/internal/config"
	"quiz-app-be/internal/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

type AwsClient struct {
	c        *s3.Client
	Bucket   string
	endpoint string
}

func Init(cfg internalConf.AwsClient) (*AwsClient, error) {
	// init s3 configuration
	s3Config, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: cfg.Endpoint},
					nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     cfg.AccessId,
				SecretAccessKey: cfg.AccessKey,
				SessionToken:    cfg.AccessToken,
			},
		}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "load aws creds error")
	}
	// init s3 client
	s := s3.NewFromConfig(s3Config,
		func(o *s3.Options) {
			o.UsePathStyle = cfg.S3ForcePathStyle
			o.EndpointOptions.DisableHTTPS = cfg.DisableSSL
		},
	)
	// ping bucket
	_, err = s.HeadBucket(
		context.Background(),
		&s3.HeadBucketInput{
			Bucket: aws.String(cfg.Bucket),
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(model.ErrNoBucket, cfg.Bucket))
	}

	return &AwsClient{
		c:        s,
		Bucket:   cfg.Bucket,
		endpoint: cfg.Endpoint,
	}, nil
}
