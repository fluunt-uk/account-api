package file

import (
	"encoding/json"
	"gitlab.com/projectreferral/account-api/internal"
	"gitlab.com/projectreferral/account-api/internal/models"
	"gitlab.com/projectreferral/account-api/lib/s3"
	"log"
	"net/http"
)

var Param = "file"

func Upload(w http.ResponseWriter, r *http.Request) {
	result, err := s3.UploadFile(r, Param)

	if err != nil {
		log.Println([]byte(err.Error()))
		http.Error(w, err.Error(), 400)
		return
	}

	if result == nil {
		s := "Upload file failed : [Result returned Nil]"
		log.Println(s)
		http.Error(w, s, 400)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func Download(w http.ResponseWriter, r *http.Request) {
	values, ok := r.URL.Query()[Param]

	if !ok || len(values[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	file, err := s3.DownloadFile(values[0])

	if err != nil {
		log.Println([]byte(err.Error()))
		http.Error(w, err.Error(), 400)
		return
	}

	if file == nil {
		s := "Download file failed : [File returned Nil]"
		log.Println(s)
		http.Error(w, s, 400)
		return
	}

	log.Printf("file: %+v", file.Name())
	w.WriteHeader(http.StatusCreated)
}

func PutEncryption(w http.ResponseWriter, r *http.Request) {
	sSEncryption := models.SSEncryption{}
	err := json.NewDecoder(r.Body).Decode(&sSEncryption)
	if !internal.HandleError(err, w) {
		result, err := s3.PutEncryption(sSEncryption.Key)
		if err != nil || result == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	}
}
