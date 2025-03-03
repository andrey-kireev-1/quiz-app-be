package aws

import (
	"bytes"
	"context"
	"fmt"
	"quiz-app-be/internal/model"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/pkg/errors"
)

const hoursInWeek = 168

func (a *AwsClient) Ping() (err error) {
	_, err = a.c.HeadBucket(
		context.Background(),
		&s3.HeadBucketInput{
			Bucket: aws.String(a.Bucket),
		})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf(model.ErrNoBucket, a.Bucket))
	}
	return
}

func (a *AwsClient) Upload(src []byte, key, contentType string) error {
	_, err := a.c.PutObject(
		context.Background(),
		&s3.PutObjectInput{
			Key:         aws.String(key),
			Bucket:      aws.String(a.Bucket),
			Body:        bytes.NewReader(src),
			ACL:         types.ObjectCannedACLPublicRead,
			ContentType: aws.String(contentType),
		})
	return err
}

func (a *AwsClient) GetDownloadUrlByName(key string) (string, error) {
	presignClient := s3.NewPresignClient(a.c)
	objRqs, err := presignClient.PresignGetObject(
		context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(a.Bucket),
			Key:    aws.String(key),
		},
		s3.WithPresignExpires(time.Hour*hoursInWeek),
	)
	return objRqs.URL, err
}

func (a *AwsClient) ListObjects(startDate time.Time, endDate time.Time) ([]string, error) {
	var objectsList []string
	var continuationToken *string
	for {
		objRqs, err := a.c.ListObjectsV2(
			context.Background(),
			&s3.ListObjectsV2Input{
				Bucket:            aws.String(a.Bucket),
				ContinuationToken: continuationToken,
			},
		)
		if err != nil {
			return nil, err
		}
		for _, file := range objRqs.Contents {
			if file.LastModified.After(startDate) && file.LastModified.Before(endDate) {
				objectsList = append(objectsList, *file.Key)
			}
		}
		if !aws.ToBool(objRqs.IsTruncated) {
			break
		}
		continuationToken = objRqs.NextContinuationToken
	}
	return objectsList, nil
}

func (a *AwsClient) Delete(key string) error {
	_, err := a.c.DeleteObject(
		context.Background(),
		&s3.DeleteObjectInput{
			Key:    aws.String(key),
			Bucket: aws.String(a.Bucket),
		},
	)
	return err
}
