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
	DC		*dynamodb.Wrapper
}
//implement only the necessary methods for each repository
//available to be consumed by the API
type AccountAdvertBuilder interface{
	GetAllAdverts(w http.ResponseWriter, r *http.Request)
}
//interface with the implemented methods will be injected in this variable
var AccountAdvert AccountAdvertBuilder

//get all the adverts for a specific account
//token validated
func (c *AccountAdvertWrapper) GetAllAdverts(w http.ResponseWriter, r *http.Request) {
	var u = models.User{}

	//email parsed from the jwt
	email := security.GetClaimsOfJWT().Subject
	result, err := 	c.DC.GetItem(email)

	if !internal.HandleError(err, w) {

		dynamodb.Unmarshal(result, &u)

		b, err := json.Marshal(u.Applications)

		if !internal.HandleError(err, w) {

			w.Write(b)
			w.WriteHeader(http.StatusOK)
		}
	}
}
