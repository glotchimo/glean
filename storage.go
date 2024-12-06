package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

var svc *s3.Client

func initStorage() error {
	// Load AWS SDK configuration
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("couldn't load default configuration: %v", err)
	}

	// Create S3 service client with Tigris configuration
	svc = s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(os.Getenv("AWS_ENDPOINT_URL_S3"))
		o.Region = "auto"
	})

	return nil
}

func listPosts() ([]string, error) {
	var posts []string

	// Create a request to list objects in the bucket
	req := &s3.ListObjectsV2Input{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
	}

	// Loop through the objects in the bucket
	paginator := s3.NewListObjectsV2Paginator(svc, req)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %v", err)
		}

		for _, item := range page.Contents {
			if strings.HasSuffix(*item.Key, ".md") {
				name := strings.TrimSuffix(*item.Key, ".md")
				posts = append(posts, name)
			}
		}
	}

	// Reverse the order to show newest first
	for i := len(posts)/2 - 1; i >= 0; i-- {
		j := len(posts) - 1 - i
		posts[i], posts[j] = posts[j], posts[i]
	}

	return posts, nil
}

func getPost(name string) ([]byte, error) {
	// Get the object from Tigris
	result, err := svc.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(name + ".md"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %v", err)
	}
	defer result.Body.Close()

	// Read the object's content
	content, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object content: %v", err)
	}

	return content, nil
}

func savePost(name string, content []byte) error {
	// Upload the file to Tigris
	_, err := svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(name + ".md"),
		Body:   bytes.NewReader(content),
	})
	if err != nil {
		return fmt.Errorf("failed to upload data: %v", err)
	}

	return nil
}
