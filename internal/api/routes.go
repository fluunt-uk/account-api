package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gitlab.com/projectreferral/account-api/configs"
	"gitlab.com/projectreferral/account-api/internal/api/account"
	account_advert "gitlab.com/projectreferral/account-api/internal/api/account-advert"
	sign_in "gitlab.com/projectreferral/account-api/internal/api/sign-in"
	"gitlab.com/projectreferral/util/pkg/security"
	"io/ioutil"
	"log"
	"net/http"
)

func SetupEndpoints() {

	_router := mux.NewRouter()

	_router.HandleFunc("/test", account.TestFunc)

	_router.HandleFunc("/upload", security.WrapHandlerWithSpecialAuth(account.UploadFile, configs.AUTH_AUTHENTICATED)).Methods("POST")

	_router.HandleFunc("/encrypt", security.WrapHandlerWithSpecialAuth(account.PutEncryption, configs.AUTH_AUTHENTICATED)).Methods("PUT")

	//token with correct register claim allowed
	_router.HandleFunc("/account", security.WrapHandlerWithSpecialAuth(account.CreateUser, configs.AUTH_REGISTER)).Methods("PUT")

	//token with correct authenticated claim allowed
	_router.HandleFunc("/account", security.WrapHandlerWithSpecialAuth(account.UpdateUser, configs.AUTH_AUTHENTICATED)).Methods("PATCH")
	_router.HandleFunc("/account", security.WrapHandlerWithSpecialAuth(account.GetUser, configs.AUTH_AUTHENTICATED)).Methods("GET")

	//no one should have access apart from super users
	//_router.HandleFunc("/account", util.WrapHandlerWithSpecialAuth(account.DeleteUser, configs.NO_ACCESS)).Methods("DELETE")

	_router.HandleFunc("/account/premium", security.WrapHandlerWithSpecialAuth(account.IsUserPremium, configs.AUTH_AUTHENTICATED)).Methods("GET")

	//token with correct sign in claim allowed
	_router.HandleFunc("/account/signin", security.WrapHandlerWithSpecialAuth(sign_in.Login, configs.AUTH_LOGIN)).Methods("POST")

	//token verification happening under the function
	_router.HandleFunc("/account/verify", security.WrapHandlerWithSpecialAuth(account.VerifyEmail, "")).Methods("POST")
	_router.HandleFunc("/account/verify/resend", security.WrapHandlerWithSpecialAuth(account.ResendVerification, "")).Methods("POST")

	//user must be authenticated before access this endpoint
	_router.HandleFunc("/account/advert", security.WrapHandlerWithSpecialAuth(account_advert.GetRefereeAdsPosted, configs.AUTH_AUTHENTICATED)).Methods("GET")
	_router.HandleFunc("/account/applications", security.WrapHandlerWithSpecialAuth(account_advert.GetJobApplications, configs.AUTH_AUTHENTICATED)).Methods("GET")
	_router.HandleFunc("/account/advert/applicants", security.WrapHandlerWithSpecialAuth(account_advert.GetAdApplicants, configs.AUTH_AUTHENTICATED)).Methods("GET")

	_router.HandleFunc("/log", displayLog).Methods("GET")

	account.Init()

	c := cors.New(cors.Options{
		AllowedMethods: []string{"POST"},
		AllowedOrigins: []string{"*"},
		AllowCredentials: true,
		AllowedHeaders: []string{"g-recaptcha-response", "Authorization", "Content-Type","Origin","Accept", "Accept-Encoding", "Accept-Language", "Host", "Connection", "Referer", "Sec-Fetch-Mode", "User-Agent", "Access-Control-Request-Headers", "Access-Control-Request-Method: "},
		OptionsPassthrough: true,
	})

	handler := c.Handler(_router)

	log.Fatal(http.ListenAndServe(configs.PORT, handler))
}

func displayLog(w http.ResponseWriter, r *http.Request){
	b, err := ioutil.ReadFile(configs.LOG_PATH)
	
	if err != nil {
			fmt.Println(err.Error()) //output to main
		w.WriteHeader(http.StatusInternalServerError)
	}else{
		w.Write(b)
	}
}