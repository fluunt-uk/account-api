package repo_builder

import (
	"encoding/json"
	"gitlab.com/projectreferral/account-api/internal"
	"gitlab.com/projectreferral/account-api/internal/models"
	"gitlab.com/projectreferral/util/pkg/dynamodb"
	"gitlab.com/projectreferral/util/pkg/security"
	"net/http"
)
type AccountAdvertWrapper struct {
	//dynamo client
	DC    *dynamodb.Wrapper
}
//implement only the necessary methods for each repository
//available to be consumed by the API
type AccountAdvertBuilder interface{
	GetRefereeAdsPosted(w http.ResponseWriter, r *http.Request)
	GetJobApplications(w http.ResponseWriter, r *http.Request)
	GetAdApplicants(w http.ResponseWriter, r *http.Request)
}
//interface with the implemented methods will be injected in this variable
var AccountAdvert AccountAdvertBuilder

//get all the adverts for a specific account
//token validated
func (c *AccountAdvertWrapper) GetRefereeAdsPosted(w http.ResponseWriter, r *http.Request) {
	var u = models.User{}

	//email parsed from the jwt
	email := security.GetClaimsOfJWT().Subject
	result, err :=     c.DC.GetItem(email)

	if !internal.HandleError(err, w) {

		dynamodb.Unmarshal(result, &u)

		b, err := json.Marshal(u.AdsPosted)

		if !internal.HandleError(err, w) {

			w.Write(b)
			w.WriteHeader(http.StatusOK)
		}
	}
}

// Get all applications
func (c *AccountAdvertWrapper) GetJobApplications(w http.ResponseWriter, r *http.Request) {
	var u = models.User{}

	//email parsed from the jwt
	email := security.GetClaimsOfJWT().Subject
	result, err :=     c.DC.GetItem(email)

	if !internal.HandleError(err, w) {

		dynamodb.Unmarshal(result, &u)

		b, err := json.Marshal(u.Applications)

		if !internal.HandleError(err, w) {

			w.Write(b)
			w.WriteHeader(http.StatusOK)
		}
	}
}

func (c *AccountAdvertWrapper) GetAdApplicants(w http.ResponseWriter, r *http.Request) {

	var ad models.Advert
	// Get ad details from when user clicks on show applicants
	errDecode := dynamodb.DecodeToMap(r.Body, &ad)

	if internal.HandleError(errDecode, w) {

		b, err := json.Marshal(ad.Applicants)

		if !internal.HandleError(err, w) {
			w.Write(b)
			w.WriteHeader(http.StatusOK)
		}
	}
	w.WriteHeader(http.StatusBadRequest)
}