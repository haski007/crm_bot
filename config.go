package main

import "os"

var (
	// TOKEN given by BotFather to use telegram API CRM_BOT_TOKEN
	TOKEN = os.Getenv("CRM_BOT_TOKEN")
	// SECRET_PASSWORD to register new admin SECRET_CRM_BOT_PASSWORD
	SECRET_PASSWORD = os.Getenv("SECRET_CRM_BOT_PASSWORD")
)

const (

	// MongoUsername ...
	MongoUsername = "haski0071"
	// MongoPassword ...
	MongoPassword = "Haski12345"
	// MongoHostname ...
	MongoHostname = "172.18.0.2"
	// MongoPort ...
	MongoPort = "27017"
)
