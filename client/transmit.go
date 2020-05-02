package client

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// TransferBatch transfer an image batch to ceph S3
func TransferBatch(path string) string {
	var (
		s3Key string
	)
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
	re := regexp.MustCompile("image.*\\.zip")
	match := re.FindSubmatch([]byte(path))
	if len(match) > 0 {
		s3Key = string(match[0])
	}

	input := &s3.PutObjectInput{
		Body:               bytes.NewReader(buffer),
		Bucket:             aws.String("test-bucket"),
		Key:                aws.String(s3Key),
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
	}
	fmt.Printf("S3 Put Object Result : %s \n", result)
	return s3Key
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

// ListObject list all objects in a bucket
func ListObject(bucket string) {
	svc := s3.New(session.New(s3Config))
	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(10),
	}

	result, err := svc.ListObjectsV2(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				fmt.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}

	fmt.Println(result)
}
