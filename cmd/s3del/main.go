package main

import (
	"log"
	"os"
	"flag"
	//"bufio"
	//"strings"
	"fmt"

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
	//delReq := &s3.DeleteBucketInput{
	//	Bucket: flgBucket,
	//}

	//reader := bufio.NewReader(os.Stdin)
	//log.Printf("Are you sure you want to permanently delete the bucket '%s' and all its contents? THIS IS NOT REVERSABLE [yes/no]", *flgBucket)
	//confirmation, _ := reader.ReadString('\n')
	//confirmation = strings.TrimRight(confirmation, "\n")
	//if confirmation != "yes" {
	//	os.Exit(1)
	//}

	// Delete all objects in a bucket
	listObjInput := &s3.ListObjectsV2Input{
		Bucket: flgBucket,
		MaxKeys: aws.Int64(5),
	}
	pageNum := 0
	err = svc.ListObjectsV2Pages(listObjInput,
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			pageNum++
			var  obj *s3.Object
			for _, obj = range page.Contents {
				fmt.Printf("OBJ: %s", *obj.Key)
				svc.DeleteObject(&s3.DeleteObjectInput{
					Bucket: flgBucket,
					Key: aws.String(obj.Key),
				})
			}
			return !page.IsTruncated()
		})

	//output, err := svc.DeleteBucket(delReq)
	//if err != nil {
	//	log.Printf("ERROR: \n\n %s", err.Error())
	//}
	//
	//log.Printf("out: %s", output)
}