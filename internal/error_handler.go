package internal

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"gitlab.com/projectreferral/util/pkg/dynamodb"
	"net/http"
)

//Custom made error
func DynamoDbError(err error, w http.ResponseWriter) bool {
	if err != nil {
		//we can return the error with specific response code and reason
		e, isCustom := err.(*dynamodb.ErrorString)

		if isCustom {
			http.Error(w, e.Reason, e.Code)
			return true
		}

		//default error
		http.Error(w, err.Error(), 400)
		return true
	}
	return false

}

func AWSError(err error, w http.ResponseWriter) bool{

	if err != nil {
		//we can return the error with specific response code and reason
		e, isCustom := err.(awserr.RequestFailure)

		if isCustom {
			http.Error(w, e.Message(), e.StatusCode())
			return true
		}

		//default error
		http.Error(w, err.Error(), 400)
		return true
	}
	return false
}

func DefaultError(err error, w http.ResponseWriter) bool{

	if err != nil {
		//default error
		http.Error(w, err.Error(), 400)
		return true
	}

	return false
}
