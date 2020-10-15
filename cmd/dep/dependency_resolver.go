package dep

import (
	"gitlab.com/projectreferral/account-api/lib/dynamodb/repo-builder"
	"gitlab.com/projectreferral/account-api/lib/rabbitmq"
	s3 "gitlab.com/projectreferral/account-api/lib/s3"
	rabbit "gitlab.com/projectreferral/util/client/rabbitmq"
	utils3 "gitlab.com/projectreferral/util/client/s3"

	"gitlab.com/projectreferral/util/pkg/dynamodb"
	"log"
)

//methods that are implemented on util
//and will be used
type ConfigBuilder interface {
	LoadEnvConfigs()
	LoadDynamoDBConfigs() *dynamodb.Wrapper
	LoadRabbitMQConfigs() *rabbit.DefaultQueueClient
	LoadS3BucketConfigs() *utils3.DefaultBucketClient
}

//internal specific configs are loaded at runtime
//takes in a object(implemented interface) of type ServiceConfigs
func Inject(builder ConfigBuilder) {

	//load the env into the object
	builder.LoadEnvConfigs()

	//setup dynamo library
	//TODO:shall the dynamo configs injected here? or in the main?
	dynamoClient := builder.LoadDynamoDBConfigs()
	//connect to the instance
	log.Println("Connecting to dynamo client")
	dynamoClient.DefaultConnect()

	//dependency injection to our resource
	//we inject the dynamo client
	//shared client, therefore shared in between all the repos
	LoadSignInRepo(&repo_builder.SignInWrapper{
		DC: dynamoClient,
	})

	LoadAccountRepo(&repo_builder.AccountWrapper{
		DC: dynamoClient,
	})

	LoadAccountAdvertRepo(&repo_builder.AccountAdvertWrapper{
		DC: dynamoClient,
	})

	//dependency injection to our resource
	//we inject the rabbitMQ client
	rabbitMQClient := builder.LoadRabbitMQConfigs()
	S3BucketClient := builder.LoadS3BucketConfigs()

	LoadRabbitMQClient(rabbitMQClient)
	LoadS3BucketClient(S3BucketClient)
}

//variable injected with the interface methods
func LoadAccountRepo(r repo_builder.AccountBuilder) {
	log.Println("Injecting Account repo")
	repo_builder.Account = r
}

//variable injected with the interface methods
func LoadAccountAdvertRepo(r repo_builder.AccountAdvertBuilder) {
	log.Println("Injecting Account Advert Repo")
	repo_builder.AccountAdvert = r
}

//variable injected with the interface methods
func LoadSignInRepo(r repo_builder.SignInBuilder) {
	log.Println("Injecting SignIn Repo")
	repo_builder.SignIn = r
}

func LoadRabbitMQClient(c rabbit.QueueClient) {
	log.Println("Injecting RabbitMQ Client")
	rabbitmq.Client = c
}

func LoadS3BucketClient(c utils3.Client) {
	log.Println("Injecting S3 Bucket Client")
	//injects the key and creates an instance of the s3 client
	c.Init()
	s3.C = c
}
