package main

import (
	"VEDA95/open_board/api/internal/config"
	"VEDA95/open_board/api/internal/db"
	appLogger "VEDA95/open_board/api/internal/log"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	if err := config.LoadEnvConfigs("./env"); err != nil {
		log.Fatal(err)
	}

	logger, err := appLogger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}

	adminUsername := flag.String("username", "", "Override admin username")
	adminEmail := flag.String("email", "", "Override admin email")
	adminPassword := flag.String("password", "", "Override admin password")

	flag.Parse()

	if len(*adminUsername) > 0 {
		os.Setenv("INITIAL_USER_USERNAME", *adminUsername)
	}

	if len(*adminEmail) > 0 {
		os.Setenv("INITIAL_USER_EMAIL", *adminEmail)
	}

	if len(*adminPassword) > 0 {
		os.Setenv("INITIAL_USER_PASSWORD", *adminPassword)
	}

	const errorMessage = "An error has occurred during the seeding process... Please view the logs for more information"
	envType := os.Getenv("ENV_TYPE")
	if len(envType) == 0 {
		envType = "development"
	}

	dbInstance, err := db.NewDB(logger)
	if err != nil {
		if envType == "production" {
			fmt.Println(errorMessage)
		}

		logger.Fatal().Err(err).Msg("An error occurred when setting up the database connection")
	}

	seeder := db.NewSeeder(dbInstance)

	if err := seeder.SeedPermissions(); err != nil {
		if envType == "production" {
			fmt.Println(errorMessage)
		}

		logger.Fatal().Err(err).Msg("An error occurred when seeding permissions")
	}

	if err := seeder.SeedRoles(); err != nil {
		if envType == "production" {
			fmt.Println(errorMessage)
		}

		logger.Fatal().Err(err).Msg("An error occurred when seeding roles")
	}

	if err := seeder.SeedInitialUser(); err != nil {
		if envType == "production" {
			fmt.Println(errorMessage)
		}

		logger.Fatal().Err(err).Msg("An error occurred when seeding the initial user")
	}

	fmt.Println("The database has been seeded successfully!")
}
