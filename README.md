# Metal Toad AWS Tools
Various Tools that use the AWS Go SDK


## Commands
1. `s3del` - Removes all objects and versioned objects from a bucket and deletes the bucket.
    

### Install
1. Install Go >= v1.7  
1. `go get -u github.com/metaltoad/aws-tools`


### Develop
1. Install Go >= v1.7
    `https://golang.org/doc/install`  
1. Clone repo into $GOPATH  
    `git clone git@github.com:metaltoad $GOPATH/src/github.com/metaltoad/aws-tools`
1. Get the AWS SDK
    `go get -u github.com/aws/aws-sdk-go`
1. Add new commands to cmd directory