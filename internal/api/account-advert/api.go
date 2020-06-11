package account_advert

import (
	"gitlab.com/projectreferral/account-api/lib/dynamodb/repo-builder"
	"net/http"
)

//get the email from the jwt
//stored in the subject claim
func GetRefereeAdsPosted(w http.ResponseWriter, r *http.Request) {
	repo_builder.AccountAdvert.GetRefereeAdsPosted(w ,r)
}

func GetJobApplications(w http.ResponseWriter, r *http.Request) {
	repo_builder.AccountAdvert.GetJobApplications(w ,r)
}

func GetAdApplicants(w http.ResponseWriter, r *http.Request) {
	repo_builder.AccountAdvert.GetAdApplicants(w ,r)
}