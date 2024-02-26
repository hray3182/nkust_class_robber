package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func SetUpEnv() {
	if os.Getenv("USERNAME") == "" {
		rootPath, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		//set env with .env file
		godotenv.Load(rootPath + "/.env")
	}
}
