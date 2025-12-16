package main

import (
	"VEDA95/open_board/api/internal/config"
	"VEDA95/open_board/api/internal/db"
	"VEDA95/open_board/api/internal/db/repository"
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

	dbInstance, err := db.NewDB()
	if err != nil {
		if envType == "production" {
			fmt.Println(errorMessage)
		}

		logger.Fatal().Err(err).Msg("An error occurred when setting up the database connection")
	}

	seeder := db.NewSeeder(dbInstance)
	generalSettingsRepo := repository.NewGeneralSettingsRepository(dbInstance)
	authSettingsRepo := repository.NewAuthSettingsRepository(dbInstance)
	emailSettingsRepo := repository.NewEmailSettingsRepository(dbInstance)

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

	generalSettings, err := generalSettingsRepo.Find()
	if err != nil {
		logger.Fatal().Err(err).Msg("An error occurred when fetching the general settings")
	}

	if generalSettings == nil {
		if err := generalSettingsRepo.Create(); err != nil {
			logger.Fatal().Err(err).Msg("An error occurred when creating the table for general settings")
		}
	}

	authSettings, err := authSettingsRepo.Find()
	if err != nil {
		logger.Fatal().Err(err).Msg("An error occurred when fetching the auth settings")
	}

	if authSettings == nil {
		if err := authSettingsRepo.Create(); err != nil {
			logger.Fatal().Err(err).Msg("An error occurred when creating the table for auth settings")
		}
	}

	emailSettings, err := emailSettingsRepo.Find()
	if err != nil {
		logger.Fatal().Err(err).Msg("An error occurred when fetching the email settings")
	}

	if emailSettings == nil {
		if err := emailSettingsRepo.Create(); err != nil {
			logger.Fatal().Err(err).Msg("An error occurred when creating the table for email settings")
		}
	}

	fmt.Println("The database has been seeded successfully!")
}
