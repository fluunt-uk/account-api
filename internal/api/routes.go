package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gitlab.com/projectreferral/account-api/configs"
	"gitlab.com/projectreferral/account-api/internal/api/account"
	account_advert "gitlab.com/projectreferral/account-api/internal/api/account-advert"
	"gitlab.com/projectreferral/account-api/internal/api/file"
	sign_in "gitlab.com/projectreferral/account-api/internal/api/sign-in"
	"gitlab.com/projectreferral/util/pkg/security"
	"io/ioutil"
	"log"
	"net/http"
)

func SetupEndpoints() {

	_router := mux.NewRouter()

	_router.HandleFunc("/test", account.TestFunc)

	_router.HandleFunc("/upload", security.WrapHandlerWithSpecialAuth(file.Upload, configs.AUTH_AUTHENTICATED)).Methods("POST", "OPTIONS")
	_router.HandleFunc("/download", security.WrapHandlerWithSpecialAuth(file.Download, configs.AUTH_AUTHENTICATED)).Methods("GET", "OPTIONS")
	_router.HandleFunc("/encrypt", security.WrapHandlerWithSpecialAuth(file.PutEncryption, configs.AUTH_AUTHENTICATED)).Methods("PUT", "OPTIONS")

	//token with correct register claim allowed
	_router.HandleFunc("/account", security.WrapHandlerWithSpecialAuth(account.CreateUser, configs.AUTH_REGISTER)).Methods("PUT", "OPTIONS")

	//token with correct authenticated claim allowed
	_router.HandleFunc("/account", security.WrapHandlerWithSpecialAuth(account.UpdateUser, configs.AUTH_AUTHENTICATED)).Methods("PATCH", "OPTIONS")
	_router.HandleFunc("/account", security.WrapHandlerWithSpecialAuth(account.GetUser, configs.AUTH_AUTHENTICATED)).Methods("GET", "OPTIONS")

	//no one should have access apart from super users
	//_router.HandleFunc("/account", util.WrapHandlerWithSpecialAuth(account.DeleteUser, configs.NO_ACCESS)).Methods("DELETE")

	_router.HandleFunc("/account/premium", security.WrapHandlerWithSpecialAuth(account.IsUserPremium, configs.AUTH_AUTHENTICATED)).Methods("GET", "OPTIONS")

	//token with correct sign in claim allowed
	_router.HandleFunc("/account/signin", security.WrapHandlerWithSpecialAuth(sign_in.Login, configs.AUTH_LOGIN)).Methods("POST", "OPTIONS")

	//token verification happening under the function
	_router.HandleFunc("/account/verify", security.WrapHandlerWithSpecialAuth(account.VerifyEmail, "")).Methods("POST", "OPTIONS")
	_router.HandleFunc("/account/verify/resend", security.WrapHandlerWithSpecialAuth(account.ResendVerification, "")).Methods("POST", "OPTIONS")

	//user must be authenticated before access this endpoint
	_router.HandleFunc("/account/advert", security.WrapHandlerWithSpecialAuth(account_advert.GetRefereeAdsPosted, configs.AUTH_AUTHENTICATED)).Methods("GET", "OPTIONS")
	_router.HandleFunc("/account/applications", security.WrapHandlerWithSpecialAuth(account_advert.GetJobApplications, configs.AUTH_AUTHENTICATED)).Methods("GET", "OPTIONS")
	_router.HandleFunc("/account/advert/applicants", security.WrapHandlerWithSpecialAuth(account_advert.GetAdApplicants, configs.AUTH_AUTHENTICATED)).Methods("GET", "OPTIONS")

	_router.HandleFunc("/log", displayLog).Methods("GET", "OPTIONS")

	c := cors.New(cors.Options{
		AllowedMethods:     []string{"POST", "PUT", "GET", "PATCH"},
		AllowedOrigins:     []string{"*"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"g-recaptcha-response", "Authorization", "Content-Type", "Origin", "Accept", "Accept-Encoding", "Accept-Language", "Host", "Connection", "Referer", "Sec-Fetch-Mode", "User-Agent", "Access-Control-Request-Headers", "Access-Control-Request-Method"},
		OptionsPassthrough: true,
	})

	handler := c.Handler(_router)

	log.Fatal(http.ListenAndServe(configs.PORT, handler))
}

func displayLog(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile(configs.LOG_PATH)

	if err != nil {
		fmt.Println(err.Error()) //output to main
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(b)
	}
}
