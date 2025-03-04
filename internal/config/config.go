package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func LoadEnvConfigs(root string) error {
	if len(root) == 0 {
		root = "."
	}

	envType := os.Getenv("ENV_TYPE")
	fileNames := []string{fmt.Sprintf("%s/.env", root)}

	if envType == "production" {
		fileNames = append(fileNames, fmt.Sprintf("%s/.env.production", root))
	} else {
		fileNames = append(fileNames, fmt.Sprintf("%s/.env.development", root))
	}

	if err := godotenv.Load(fileNames...); err != nil {
		return err
	}

	return nil
}
