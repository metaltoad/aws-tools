package main

import (
	"log"
	"os"
	"flag"
	"bufio"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

var flgBucket *string

func init() {
	usage := "supply a bucket to delete"
	flgBucket = flag.String("bucket", "", usage)
	flgBucket = flag.String("b", "", usage + " (shorthand)")
}

func main() {
	startTime := time.Now()
	flag.Parse()
	if *flgBucket == "" {
		log.Fatal("Please supply an s3 bucket.")
	}

	svc := initAWSSvc()

	if confirmDelete() != "yes" {
		os.Exit(1)
	}

	err := deleteObjects(svc)
	if err != nil {
		log.Fatalf("### Error deleteing objects:\n %s", err.Error())
	}

	err = deleteVersions(svc)
	if err != nil {
		log.Fatalf("### Error deleteing versions: \n %s", err.Error())
	}

	output, err := deleteBucket(svc)
	if err != nil {
		log.Printf("### Error Deleting Bucket: \n %s", err.Error())
	}

	log.Printf("### Bucket Succesfully Deleted: %s", output.GoString())
	elapsed := time.Since(startTime)
	log.Printf("Time Taken: %s \n\n", elapsed.Minutes())
}

//initAWSService returns an initialized s3 client
func initAWSSvc() *s3.S3 {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err.Error())
	}

	sdkLoadConfig := os.Getenv("AWS_SDK_LOAD_CONFIG")
	if sdkLoadConfig != "true" {
		log.Fatal(`Env var "AWS_SDK_LOAD_CONFIG" needs to be true to read credentials.\n\n Run "export AWS_SDK_LOAD_CONFIG=true" to fix this. Aborting run.`)
	}

	return s3.New(sess)
}

//confirmDelete prompts the user to confirm they want to issue the s3 bucket delete operation
func confirmDelete() string {
	reader := bufio.NewReader(os.Stdin)
	log.Printf("Are you sure you want to permanently delete the bucket '%s' and all its contents? THIS IS NOT REVERSIBLE! [yes/no]", *flgBucket)
	confirmation, _ := reader.ReadString('\n')
	return strings.TrimRight(confirmation, "\n")
}

//deleteObjects removes all objects from s3 bucket
func deleteObjects(svc *s3.S3) error {
	// Delete all objects in a bucket
	listObjInput := &s3.ListObjectsV2Input{
		Bucket: flgBucket,
		MaxKeys: aws.Int64(5),
	}
	log.Print("########## DELETEING ALL OBJECTS IN BUCKET ##########")
	err := svc.ListObjectsV2Pages(listObjInput,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			var obj *s3.Object
			for _, obj = range page.Contents {
				log.Printf("Deleting Object: %s \n", *obj.Key)
				svc.DeleteObject(&s3.DeleteObjectInput{
					Bucket: flgBucket,
					Key: obj.Key,
				})
			}
			return *page.IsTruncated
		})
	return err
}

//deleteVersions removes all versioned objects from an s3 bucket
func deleteVersions(svc *s3.S3) error {
	log.Print("########## DELETE ALL VERSIONED OBJECTS IN BUCKET ##########")
	// Delete all versioned objects in bucket
	listObjVerInput := &s3.ListObjectVersionsInput{
		Bucket: flgBucket,
		MaxKeys: aws.Int64(5),
	}
	err := svc.ListObjectVersionsPages(listObjVerInput,
		func(page *s3.ListObjectVersionsOutput, lastPage bool) bool {
			var obj *s3.ObjectVersion
			for _, obj = range page.Versions {
				log.Printf("Deleting Version: %s \n", *obj.VersionId)
				svc.DeleteObject(&s3.DeleteObjectInput{
					Bucket: flgBucket,
					Key: obj.Key,
					VersionId: obj.VersionId,
				})
			}
			return *page.IsTruncated
		})
	return err
}

//deleteBucket removes an empty bucket
func deleteBucket(svc *s3.S3) (*s3.DeleteBucketOutput, error) {
	delReq := &s3.DeleteBucketInput{
		Bucket: flgBucket,
	}
	return svc.DeleteBucket(delReq)
}



