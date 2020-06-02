package sign_in

import (
	"gitlab.com/projectreferral/account-api/lib/dynamodb/repo-builder"
	"net/http"
)

//credentials extract from the body
//query the db with the email
//if this exists, get the pw hash and compare
func Login(w http.ResponseWriter, r *http.Request) {
	repo_builder.SignIn.Login(w, r)
}
