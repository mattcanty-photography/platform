package main

import (
	"encoding/json"
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v3/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func createPhotosResources(ctx *pulumi.Context) (s3.Bucket, error) {
	bucket, err := s3.NewBucket(
		ctx,
		awsNamePrintf(ctx, "%s", "photos"),
		&s3.BucketArgs{})

	tmpJSON := bucket.Arn.ApplyT(func(arn string) (string, error) {
		policyJSON, err := json.Marshal(map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []map[string]interface{}{
				{
					"Effect":    "Allow",
					"Principal": "*",
					"Action":    []string{"s3:GetObject"},
					"Resource": []string{
						fmt.Sprintf("%s/*", arn),
					},
				},
			},
		})
		if err != nil {
			return "", err
		}
		return string(policyJSON), nil
	})
	if err != nil {
		return s3.Bucket{}, err
	}

	s3.NewBucketPolicy(
		ctx,
		awsNamePrintf(ctx, "%s", "photos"),
		&s3.BucketPolicyArgs{
			Bucket: bucket.ID(),
			Policy: tmpJSON,
		},
	)

	return *bucket, err
}