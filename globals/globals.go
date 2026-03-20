package globals

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadenv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: Unable to find .env")
	}
}

type VarsStruct struct {
	MONGO_URI     string
	DATABASE_NAME string
	PORT          string
	ACCESS_TOKEN_SECRET string
	REFRESH_TOKEN_SECRET string
	ACCESS_TOKEN_EXPIRY_MINUTES int
	REFRESH_TOKEN_EXPIRY_MINUTES int
}

var Vars VarsStruct

func init() {
	loadenv()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	Vars = VarsStruct{
		MONGO_URI:     os.Getenv("MONGO_URI"),
		DATABASE_NAME: os.Getenv("DATABASE_NAME"),
		PORT:          port,
		ACCESS_TOKEN_SECRET: os.Getenv("ACCESS_TOKEN_SECRET"),
		REFRESH_TOKEN_SECRET: os.Getenv("REFRESH_TOKEN_SECRET"),
		ACCESS_TOKEN_EXPIRY_MINUTES: 30,
		REFRESH_TOKEN_EXPIRY_MINUTES: 60 * 24 * 7,
	}
}
