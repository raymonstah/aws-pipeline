package awsutils

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// UploadLambdas find any zip files in the given dir and upload them to the lambdasBucket
func UploadLambdas(ctx context.Context, s3API s3iface.S3API, dir string, lambdasBucket string) error {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("unable to read target directory: %v: %w", dir, err)
	}

	err = createBucketIfNotExists(ctx, s3API, lambdasBucket)
	if err != nil {
		return fmt.Errorf("unable to create bucket: %w", err)
	}

	for _, fInfo := range fileInfos {
		if strings.HasSuffix(fInfo.Name(), "zip") {
			// ok
			f, err := os.Open(dir + fInfo.Name())
			defer f.Close()
			if err != nil {
				return fmt.Errorf("unable to open file %v: %w", fInfo.Name(), err)
			}
			key := fInfo.Name()
			fmt.Printf("uploading to s3://%v/%v\n", lambdasBucket, key)
			output, err := s3API.PutObjectWithContext(ctx, &s3.PutObjectInput{
				Bucket: aws.String(lambdasBucket),
				Key:    aws.String(key),
				Body:   f,
			})
			if err != nil {
				return fmt.Errorf("unable to upload file to s3://%v/%v: %w", lambdasBucket, key, err)
			}

			if output.VersionId == nil {
				return fmt.Errorf("lambdasBucket %v must support versioning", lambdasBucket)
			}
		}
	}
	return nil
}

// creates a s3 bucket if it doesn't exist
func createBucketIfNotExists(ctx aws.Context, s3API s3iface.S3API, bucketName string) error {
	_, err := s3API.HeadBucketWithContext(ctx, &s3.HeadBucketInput{Bucket: aws.String(bucketName)})
	if err != nil {

		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
			case "NotFound":
				if err := createVersionedBucket(ctx, s3API, bucketName); err != nil {
					return err
				}
				return nil
			default:
				fmt.Println("unexpected error code:", aerr.Code(), aerr.Error())
				return err
			}
		}
		return fmt.Errorf("unable to head bucket %v: %w", bucketName, err)
	}

	fmt.Printf("bucket %v already exists\n", bucketName)
	return nil
}

// create a versioned bucket in s3 us-west-2
func createVersionedBucket(ctx aws.Context, s3API s3iface.S3API, bucketName string) error {

	_, err := s3API.CreateBucketWithContext(ctx,
		&s3.CreateBucketInput{
			Bucket:                    aws.String(bucketName),
			CreateBucketConfiguration: &s3.CreateBucketConfiguration{LocationConstraint: aws.String("us-west-2")},
		})
	if err != nil {
		return fmt.Errorf("unable to create bucket %v: %w", bucketName, err)
	}

	_, err = s3API.PutBucketVersioningWithContext(ctx, &s3.PutBucketVersioningInput{
		Bucket: aws.String(bucketName),
		VersioningConfiguration: &s3.VersioningConfiguration{
			Status: aws.String("Enabled"),
		},
	})
	if err != nil {
		return fmt.Errorf("unable to enable versioning for bucket %v: %w", bucketName, err)
	}

	fmt.Printf("bucket %v created\n", bucketName)
	return nil
}
