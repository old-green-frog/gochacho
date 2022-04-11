package conf

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Host           string
	Port           uint32
	DatabaseString string
}

func New() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error: not find .env file in project root")
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}

	return &Config{
		Host:           os.Getenv("HOST"),
		Port:           uint32(port),
		DatabaseString: os.Getenv("CONNECT"),
	}
}
