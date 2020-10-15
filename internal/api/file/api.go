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

	if err != nil || result == nil {
		if err != nil {
			log.Println([]byte(err.Error()))
		}
		w.WriteHeader(http.StatusBadRequest)
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
	if err != nil || file == nil {
		if err != nil {
			log.Println([]byte(err.Error()))
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s, _ := file.Stat()

	log.Printf("file: %+v", s)
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
