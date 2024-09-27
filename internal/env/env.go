package env

import (
	"github.com/joho/godotenv"
)

func LoadEnvs() error {
	err := godotenv.Load()
	return err
}
