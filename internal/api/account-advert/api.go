package account_advert

import (
	"gitlab.com/projectreferral/account-api/lib/dynamodb/repo-builder"
	"net/http"
)

//get the email from the jwt
//stored in the subject claim
func GetAllAdverts(w http.ResponseWriter, r *http.Request) {
	repo_builder.AccountAdvert.GetAllAdverts(w ,r)
}

func GetAllApplications(w http.ResponseWriter, r *http.Request) {
	repo_builder.AccountAdvert.GetAllApplications(w ,r)
}