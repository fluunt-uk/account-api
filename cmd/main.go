package main

import (
	"gitlab.com/projectreferral/account-api/cmd/dep"
	"gitlab.com/projectreferral/account-api/configs"
	"gitlab.com/projectreferral/account-api/internal/api"
	"gitlab.com/projectreferral/account-api/internal/models"
	"gitlab.com/projectreferral/util/util"
	"log"
	"os"
)

func main() {

	f, err := os.OpenFile("logs/accountAPI_log.txt", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	//gets all the necessary configs into our object
	//completes connections
	//assigns connections to repos
	dep.Inject(&util.ServiceConfigs{
		Environment: os.Getenv("ENV"),
		Region:       configs.EU_WEST_2,
		Table:        configs.TABLE_NAME,
		SearchParam:  configs.UNIQUE_IDENTIFIER,
		GenericModel: models.User{},
		BrokerUrl:    configs.QAPI_URL,
		Port:		  configs.PORT,
	})

	api.SetupEndpoints()
}

