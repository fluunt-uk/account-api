package api

import (
	"encoding/json"
	"gitlab.com/projectreferral/account-api/configs"
	"gitlab.com/projectreferral/account-api/internal/models"
	request "gitlab.com/projectreferral/util/pkg/http_lib"
	"log"
	"net/http"
	"net/url"
)

func RecaptchaVerify(w *http.ResponseWriter, token *string, r *models.ReCaptcha){

	form := url.Values{}
	form.Add("response", *token)
	form.Add("secret", configs.RECAPTCHA_SECRET)

	reqVer, errReq := request.Post(configs.RECAPTCHA_VERIFY, []byte(form.Encode()), map[string]string{
		"Content-Type": "application/x-www-form-urlencoded"})

	if errReq != nil {
		log.Printf("Request to [%s] failed\n", configs.RECAPTCHA_VERIFY)
		http.Error(*w, errReq.Error(), 400)
		return
	}

	json.NewDecoder(reqVer.Body).Decode(&r)
}

