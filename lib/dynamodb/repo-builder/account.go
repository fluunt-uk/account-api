package repo_builder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gitlab.com/projectreferral/account-api/configs"
	"gitlab.com/projectreferral/account-api/internal"
	"gitlab.com/projectreferral/account-api/internal/models"
	"gitlab.com/projectreferral/account-api/lib/rabbitmq"
	"gitlab.com/projectreferral/util/pkg/dynamodb"
	"gitlab.com/projectreferral/util/pkg/security"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)


type AccountWrapper struct {
	//dynamo client
	DC		*dynamodb.Wrapper
}
//implement only the necessary methods for each repository
//available to be consumed by the API
type AccountBuilder interface{
	GetUser(http.ResponseWriter, *http.Request)
	UpdateUser(http.ResponseWriter, *http.Request)
	CreateUser(http.ResponseWriter, *http.Request)
	IsUserPremium(http.ResponseWriter, *http.Request)
	VerifyEmail(http.ResponseWriter, *http.Request)
	ResendVerification(http.ResponseWriter, *http.Request)
	UploadFile(http.ResponseWriter, *http.Request)
}
//interface with the implemented methods will be injected in this variable
var Account AccountBuilder

func (c *AccountWrapper) UploadFile(w http.ResponseWriter, r *http.Request) {

	// Limit size to 10 mb
	_ = r.ParseMultipartForm(10 * 1024 * 1024)

	file, handle, err := r.FormFile("File")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	defer file.Close()

	var buf bytes.Buffer
	io.Copy(&buf, file)

	i, er := ioutil.ReadAll(&buf)

	if er != nil {
		fmt.Println("failed to convert")
	}

	fmt.Println("File name:", handle.Filename)
	fmt.Println("File size:", handle.Size)
	fmt.Println("File type:", handle.Header.Get("Content-Type"))

	// Need to move this to configs or somewhere else
	s, err := session.NewSession(&aws.Config{Region: aws.String(configs.EU_WEST_2)})
	if err != nil {
		log.Fatal(err)
	}

	svc := s3.New(s)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		fmt.Println("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}

	// Will need an if statement here for
	s3er := AddFileToS3(s, i, handle.Filename, handle.Size)

	if s3er != nil {
		fmt.Println(s3er)
		w.WriteHeader(http.StatusBadRequest)
	}
	fmt.Println("File uploaded")
	w.WriteHeader(http.StatusCreated)
}

func AddFileToS3(s *session.Session, buffer []byte, filename string, size int64) error {

	// Open the file for use
	//file, err := os.Open(fileDir)
	//if err != nil {
	//	return err
	//}
	//defer file.Close()


	//_, err = buffer.WriteTo(w)
	//// Get file size and read the file content into a buffer
	//fileInfo, _ := file.Stat()
	//var size int64 = fileInfo.Size()
	//buffer := make([]byte, size)
	//file.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	worked, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(configs.S3_BUCKET),
		Key:                  aws.String(filename),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	fmt.Println(worked)
	return err
}

//We check for the recaptcha response and proceed
//Covert the response body into appropriate models
//Create a new user using our dynamodb adapter
//A event message it sent to the queues which are consumed by the relevant services
func (c *AccountWrapper) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User

	//TODO: reCaptcha check, 30ms average
	if r.ContentLength < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No body error"))
		return
	}
	body := r.Body

	u.AccessCode = rabbitmq.NewUUID()

	dynamoAttr, errDecode := dynamodb.DecodeToDynamoAttribute(body, &u)
	dynamodb.AddEmptyCollection(dynamoAttr, configs.ACTIVE_SUB)
	dynamodb.AddEmptyCollection(dynamoAttr, configs.APPLICATIONS)

	if !internal.HandleError(errDecode, w) {
		err := 	c.DC.CreateItem(dynamoAttr)

		if !internal.HandleError(err, w) {

			b, err := json.Marshal(u)
			if err != nil {
				fmt.Sprintf(err.Error())
			}
			//JSON format of the newly created user
			w.Write(b)
			w.WriteHeader(http.StatusOK)

			//triggers email confirmation e-mail
			go rabbitmq.BroadcastUserCreatedEvent(b)
		}
	}
}

//get the email from the jwt
//stored in the subject claim
func (c *AccountWrapper) GetUser(w http.ResponseWriter, r *http.Request) {
	var u models.User

	//email parsed from the jwt
	//email := security.GetClaimsOfJWT().Subject
	result, err := c.DC.GetItem("lunos4@gmail.com")

	if !internal.HandleError(err, w) {
		dynamodb.Unmarshal(result, &u)
		b, err := json.Marshal(&u)

		if !internal.HandleError(err, w) {

			w.Write(b)
			w.WriteHeader(http.StatusOK)
		}
	}
}

//two ways of updating a user's information
//type 1: updates a single string value for a defined field
//type 2: appends a map for a defined field(this field name must already exists)
//all parameters are set under ChangeRequest struct
func (c *AccountWrapper) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var cr models.ChangeRequest

	dynamodb.DecodeToMap(r.Body, &cr)

	email := security.GetClaimsOfJWT().Subject
	err := c.UpdateValue(email, &cr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	log.Printf("Updated account details for [%s] to [%v]",email, &cr)
	w.WriteHeader(http.StatusOK)
}

//check if the user has an active subscription
//parses email from the jwt
func (c *AccountWrapper) IsUserPremium(w http.ResponseWriter, r *http.Request) {
	//email parsed from the jwt
	email := security.GetClaimsOfJWT().Subject
	result, err := c.DC.GetItem(email)

	p := result.Item[configs.PREMIUM].BOOL

	if !internal.HandleError(err, w) && *p {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(204)
	return
}

//we parse the access_code and token from the query string
//token is validated
//we compare the access_code in the db matches the one passed in from the query string
func (c *AccountWrapper) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	queryMap := r.URL.Query()

	accessCodeKeys, ok := queryMap["access_code"]
	tokenKeys, ok := queryMap["token"]
	if !ok {
		w.Write([]byte("Url Param are missing"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	accessCodeValue := accessCodeKeys[0]
	tokenValue := tokenKeys[0]
	if len(accessCodeValue) < 1 || len(tokenValue) < 1 {
		w.Write([]byte("Url Param are missing"))
		w.WriteHeader(http.StatusBadRequest)
	}

	//validate the token expiry date
	if security.VerifyTokenWithClaim(tokenValue, configs.AUTH_VERIFY) {

		//email parsed from the jwt
		email := security.GetClaimsOfJWT().Subject
		result, err := c.DC.GetItem("lunos4@gmail.com")

		if !internal.HandleError(err, w) {
			if accessCodeValue == *result.Item["access_code"].S {

				c.UpdateValue(email, &models.ChangeRequest{Field: "verified", NewBool: true, Type: 3})
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Access code does not match"))
			return
		}
	}

	w.WriteHeader(http.StatusBadRequest)
}

//TODO:resend verification email
func (c *AccountWrapper) ResendVerification(w http.ResponseWriter, r *http.Request) {
	var u models.User
	email := security.GetClaimsOfJWT().Subject

	//new access code generated
	c.UpdateValue(email, &models.ChangeRequest{Field: "access_code", NewString: rabbitmq.NewUUID(), Type: 1})

	user, err := c.DC.GetItem("lunos4@gmail.com")


	if !internal.HandleError(err, w) {

		dynamodb.Unmarshal(user, &u)
		b, errM := json.Marshal(&u)

		if !internal.HandleError(errM, w) {

			w.Write(b)
			w.WriteHeader(http.StatusOK)

			go rabbitmq.BroadcastUserCreatedEvent(b)
		}
	}
}