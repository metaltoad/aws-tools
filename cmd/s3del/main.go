package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"flag"
	"bufio"
	"strings"
)

var flgBucket *string

func init() {
	usage := "supply a bucket to delete"
	flgBucket = flag.String("bucket", "", usage)
	flgBucket = flag.String("b", "", usage + " (shorthand)")
}
func main() {
	flag.Parse()
	if *flgBucket == "" {
		log.Fatal("Please supply an s3 bucket.")
	}

	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err.Error())
	}

	sdkLoadConfig := os.Getenv("AWS_SDK_LOAD_CONFIG")
	if sdkLoadConfig != "true" {
		log.Fatal(`Env var "AWS_SDK_LOAD_CONFIG" needs to be true to read credentials.\n\n Run "export AWS_SDK_LOAD_CONFIG=true" to fix this. Aborting run.`)
	}

	svc := s3.New(sess)
	delReq := &s3.DeleteBucketInput{
		Bucket: flgBucket,
	}

	reader := bufio.NewReader(os.Stdin)
	log.Printf("Are you sure you want to permanently delete the bucket '%s' and all its contents? THIS IS NOT REVERSABLE [yes/no]", *flgBucket)
	confirmation, _ := reader.ReadString('\n')
	log.Printf("INPUT: %s", confirmation)
	if strings.Compare(confirmation, "yes") == 0 {
		os.Exit(1)
	}

	_, err = svc.DeleteBucket(delReq)
}