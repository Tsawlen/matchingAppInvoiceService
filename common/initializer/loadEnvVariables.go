package initializer

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables(readyChannel chan bool) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	readyChannel <- true
}
