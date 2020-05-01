package client

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// TransferBatch transfer an image batch to ceph S3
func TransferBatch(path string) {
	// Open the file for use
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("err: %s \n", err.Error())
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	s3Client := s3.New(session.New(s3Config))
	// Regex gets batch name

	input := &s3.PutObjectInput{
		Body:               bytes.NewReader(buffer),
		Bucket:             aws.String("test-bucket"),
		Key:                aws.String("image_batch_1"),
		ContentLength:      aws.Int64(size),
		ContentType:        aws.String(http.DetectContentType(buffer)),
		ContentDisposition: aws.String("image-batch"),
	}
	result, err := s3Client.PutObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}
	fmt.Println(result)
}

// ListBucket lists ceph s3 bucket
func ListBucket() {
	s3Client := s3.New(session.New(s3Config))
	result, err := s3Client.ListBuckets(nil)
	if err != nil {
		fmt.Printf("Unable to list bucket, %s\n", err)
	}

	fmt.Println("Bucket: ")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s \n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
}
