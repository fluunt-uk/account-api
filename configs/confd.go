package configs

const (
	PORT = ":5001"
	LOG_PATH         = "../logs/accountAPI_log.txt"
	/************** DynamoDB configs *************/
	EU_WEST_2         = "eu-west-2"
	UNIQUE_IDENTIFIER = "email"
	PW                = "password"
	PREMIUM           = "premium"
	APPLICATIONS      = "applications"
	ACTIVE_SUB        = "active_subscription"
	TABLE_NAME        = "users"
	/*********************************************/
	/************** RabbitMQ configs *************/
	FANOUT_EXCHANGE = "accounts.fanout"
	//for dev usage outside of local network
	//QAPI_URL = "http://35.179.11.178:5004"
	QAPI_URL = "http://localhost:5004"
	/*********************************************/
	/*********** Authentication(token permissions) configs **********/
	AUTH_REGISTER      = "register_user"
	AUTH_AUTHENTICATED = "crud"
	AUTH_LOGIN         = "signin_user"
	AUTH_VERIFY        = "verify_user"
	NO_ACCESS          = "admin_gui"
	/*********************************************/
	/*************** S3 configs ******************/
	S3_BUCKET		   = "docs-s3"
	S3_ENCRYPTION_ALGORITHM    = "AES256"
	PART_SIZE                  = 10 * 1024 * 1024
	S3_DOWNLOAD_LOCATION       = "../cache/"
	/*****************************************************************/
	/*************************** Google Recaptcha(I am not a Robot) configs **************************/
	RECAPTCHA_VERIFY = "https://www.google.com/recaptcha/api/siteverify"
	RECAPTCHA_SECRET = "6LcbrKIZAAAAACoS7IHx5KZfhkk3T1tXBhcIGf6W"
	/*****************************************************************/
)

var (
	//To dial RabbitMQ
	BrokerUrl = ""
)
