package main

import (
	"VEDA95/open_board/api/internal/config"
	"VEDA95/open_board/api/internal/db"
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/validators"
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

	username := flag.String("username", "", "Username (required)")
	password := flag.String("password", "", "Password (min 8 chars)")
	confirmPassword := flag.String("confirm_password", "", "Confirm Password (must be the same as the password)")

	flag.Parse()

	const errorMessage = "An error has occurred during the during the password reset process... Please view the logs for more information"
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

	userRepo := repository.NewUserRepository(dbInstance)
	sessionRepo := repository.NewSessionRepository(dbInstance)
	validator := validators.NewValidator()
	validatorData := &validators.ResetPasswordProgramValidator{
		Username:        *username,
		Password:        *password,
		ConfirmPassword: *confirmPassword,
	}

	if errs := validator.Validate(validatorData); len(errs) > 0 {
		if envType == "production" {
			fmt.Println(errorMessage)

			for _, err := range errs {
				fmt.Printf("%s: %s\n", err.FailedField, err.ErrValue)
			}
		}

		logger.Fatal().Err(errors.CreateValidationError(errs)).Msg("An error occurred during validation")
	}

	user, err := userRepo.FindByUsername(validatorData.Username, repository.QueryOptions{
		Omit: []string{"Roles", "Sessions"},
	})
	if err != nil {
		if envType == "production" {
			fmt.Println(errorMessage)
		}

		logger.Fatal().Err(err).Msg("An error occurred while fetching the user")
	}

	if err := user.HashPassword(validatorData.Password); err != nil {
		if envType == "production" {
			fmt.Println(errorMessage)
		}

		logger.Fatal().Err(err).Msg("An error occurred while hashing the user password")
	}

	if err := userRepo.Update(user, repository.QueryOptions{Select: []string{"hashed_password"}}); err != nil {
		if envType == "production" {
			fmt.Println(errorMessage)
		}

		logger.Fatal().Err(err).Msg("An error occurred when updating user password")
	}

	if err := sessionRepo.DeleteByUserID(user.ID); err != nil {
		if envType == "production" {
			fmt.Println(errorMessage)
		}

		logger.Fatal().Err(err).Msg("An error occurred while deleting user sessions")
	}

	fmt.Printf("The password for user %s has been reset successfully!!!\n", user.Username)
}
