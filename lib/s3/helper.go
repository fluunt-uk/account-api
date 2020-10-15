package s3

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	s3bucket "gitlab.com/projectreferral/util/client/s3"
	"log"
	"net/http"
	"os"
)

var C s3bucket.Client

func UploadFile(r *http.Request, name string) (*s3manager.UploadOutput, error) {

	if name != "" || r != nil {
		return C.UploadFile(r, name)
	}

	return nil, errors.New("file name undefined")
}

func DownloadFile(name string) (*os.File, error) {

	if name != "" {
		return C.DownloadFile(name)
	}

	return nil, errors.New("file name undefined")
}

//function for putting KMS key
func PutEncryption(key string) (*s3.PutBucketEncryptionOutput, error) {

	if key != "" {
		return C.PutEncryption(key)
	}

	return nil, errors.New("key undefined")
}

//Custom made error
func HandleError(err error) bool {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Println(aerr.Error())
			}
		} else {
			log.Println(err.Error())
		}
		return true
	}
	return false
}
