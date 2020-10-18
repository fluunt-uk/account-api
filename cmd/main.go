package main

import (
	"gitlab.com/projectreferral/account-api/cmd/dep"
	"gitlab.com/projectreferral/account-api/configs"
	"gitlab.com/projectreferral/account-api/internal/api"
	"gitlab.com/projectreferral/account-api/internal/models"
	docbucket "gitlab.com/projectreferral/util/client/s3/models"
	"gitlab.com/projectreferral/util/util"
	"log"
	"os"
)

func main() {
	f, err := os.OpenFile(configs.LOG_PATH, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	//log all to the file, disabled for local debugging
	if os.Getenv("ENV") == "prod" {log.SetOutput(f)}

	//gets all the necessary configs into our object
	//completes connections
	//assigns connections to repos
	dep.Inject(&util.ServiceConfigs{
		Environment:	os.Getenv("ENV"),
		Region:       	configs.EU_WEST_2,
		Table:        	configs.TABLE_NAME,
		SearchParam:  	configs.UNIQUE_IDENTIFIER,
		GenericModel: 	models.User{},
		BrokerUrl:    	configs.QAPI_URL,
		Port:		  	configs.PORT,
		S3Config:     	&docbucket.S3Configs{
			Region:              configs.EU_WEST_2,
			Key:                 os.Getenv("S3_KEY"),
			DownloadLocation:    configs.S3_DOWNLOAD_LOCATION,
			Bucket:              configs.S3_BUCKET,
			EncryptionAlgorithm: configs.S3_ENCRYPTION_ALGORITHM,
			PartSize:            configs.PART_SIZE,
		},
	})

	api.SetupEndpoints()
}

