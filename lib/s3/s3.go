package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"gitlab.com/projectreferral/account-api/configs"
	"log"
	"net/http"
	"os"
)

var s3Key string
var s3Session *session.Session

func Init(){
	s3Key = os.Getenv(configs.S3_KEY)
	if s3Key == "" {
		log.Println("No s3 key found")
		os.Exit(1)
	}
	s3Session = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(configs.EU_WEST_2)},
	))
	if s3Session == nil {
		log.Println("Initiation failed")
		os.Exit(1)
	}
	log.Println("Session Initiated")
}

func UploadFile(r *http.Request, name string) (*s3manager.UploadOutput,error) {
	sizeErr := r.ParseMultipartForm(configs.PART_SIZE) // Limit size to part size

	if !HandleError(sizeErr)  {

		file, header, fErr := r.FormFile(name)
		if !HandleError(fErr) {

			defer file.Close()

			filename := header.Filename
			size := header.Size
			log.Printf("file name %s size %d",filename,size)

			input := &s3manager.UploadInput{
				// Bucket to be used
				Bucket: aws.String(configs.S3_BUCKET),
				// Name of the file to be saved
				Key:    aws.String(filename),
				// File body
				Body:   file,
				// Encrypt file
				SSECustomerAlgorithm: aws.String(configs.S3_ENCRYPTION_ALGORITHM),
				SSECustomerKey : aws.String(s3Key),
			}
			uploader := s3manager.NewUploader(s3Session)
			log.Println("uploader created")

			// Perform upload with multipart
			result, UErr := uploader.Upload(input, func(u *s3manager.Uploader) {
				u.PartSize = configs.PART_SIZE
				u.LeavePartsOnError = true    // Don't delete the parts if the upload fails.
			})

			if !HandleError(UErr){
				log.Printf("%+v",result)
				return result, nil
			}
		}
		return nil, fErr
	}
	return nil, sizeErr
}

func DownloadFile(name string) (*os.File,int64,error) {
	file, fErr := os.Create(configs.S3_DOWNLOAD_LOCATION)
	if !HandleError(fErr) && file == nil {
		log.Println("error creating file")
		return nil, -1, nil
	}

	defer file.Close()

	downloader := s3manager.NewDownloader(s3Session)
	size, err := downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(configs.S3_BUCKET),
		Key:    aws.String(name),
		SSECustomerAlgorithm: aws.String(configs.S3_ENCRYPTION_ALGORITHM),
		SSECustomerKey : aws.String(s3Key),

	})
	if !HandleError(err) {
		return file, size, err
	}

	return file, size, nil
}

//function for putting KMS key
func PutEncryption(key string) (*s3.PutBucketEncryptionOutput,error) {
	defEnc := &s3.ServerSideEncryptionByDefault{KMSMasterKeyID: aws.String(key), SSEAlgorithm: aws.String(configs.S3_ENCRYPTION_ALGORITHM)}
	rule := &s3.ServerSideEncryptionRule{ApplyServerSideEncryptionByDefault: defEnc}
	rules := []*s3.ServerSideEncryptionRule{rule}
	serverConfig := &s3.ServerSideEncryptionConfiguration{Rules: rules}
	input := &s3.PutBucketEncryptionInput{Bucket: aws.String(configs.S3_BUCKET), ServerSideEncryptionConfiguration: serverConfig}
	svc := s3.New(s3Session)
	result, err := svc.PutBucketEncryption(input)
	if !HandleError(err) {
		log.Println("Bucket %s now has KMS encryption by default %+v",configs.S3_BUCKET,result)
		return result, nil
	}
	return result, err
}

//Custom made error
func HandleError(err error) bool {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
				default:
					log.Println(aerr.Error())
			}
		}else {
			log.Println(err.Error())
		}
		return true
	}
	return false
}
