package main

import (
	"VEDA95/open_board/api/internal/config"
	"VEDA95/open_board/api/internal/db"
	models "VEDA95/open_board/api/internal/db/model"
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/http/validators"
	appLogger "VEDA95/open_board/api/internal/log"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"gorm.io/gorm"
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
	email := flag.String("email", "", "Email (required)")
	password := flag.String("password", "", "Password (min 8 chars)")
	confirmPassword := flag.String("confirm_password", "", "Confirm Password (must be the same as the password)")
	firstName := flag.String("first", "", "First name")
	lastName := flag.String("last", "", "Last name")
	roles := flag.String("roles", "", "Comma-separated role names")
	superuser := flag.Bool("superuser", false, "Create as superuser")
	admin := flag.Bool("admin", false, "Create as admin")

	flag.Parse()

	const errorMessage = "An error has occurred during the during the user creation process... Please view the logs for more information"
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
	validator := validators.NewValidator()
	splitRoles := strings.Split(*roles, ",")
	validatorData := &validators.CreateUserProgramValidator{
		Username:        *username,
		Email:           *email,
		Password:        *password,
		ConfirmPassword: *confirmPassword,
		SuperUser:       *superuser,
	}

	if len(*firstName) > 0 {
		validatorData.FirstName = firstName
	}

	if len(*lastName) > 0 {
		validatorData.LastName = lastName
	}

	if errs := validator.Validate(validatorData); len(errs) > 0 {
		if envType == "production" {
			fmt.Println(errorMessage)

			for _, err := range errs {
				fmt.Printf("%s: %s\n", err.FailedField, err.ErrValue)
			}
		}

		logger.Fatal().Interface("errors", errs).Msg("The following validation errors occurred")
	}

	var userCount int64
	if err := dbInstance.Model(&models.User{}).Count(&userCount).Error; err != nil {
		if envType == "production" {
			fmt.Println(errorMessage)
		}

		logger.Fatal().Err(err).Msg("An error occurred during the user creation process")
	}

	existingUserByUsername, err := userRepo.FindByUsername(validatorData.Username, repository.QueryOptions{
		Select: []string{"id"},
	})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		if envType == "production" {
			fmt.Println(errorMessage)
		}

		logger.Fatal().Err(err).Msg("An error occurred when checking if user already exists")
	}

	if existingUserByUsername != nil {
		err := errors.New("A user with the provided username already exists")

		if envType == "production" {
			fmt.Print(err.Error())
		}

		logger.Fatal().Err(err).Msg("An error occurred when checking if the user already exists")
	}

	existingUserByEmail, err := userRepo.FindByEmail(validatorData.Email, repository.QueryOptions{
		Select: []string{"id"},
	})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		if envType == "production" {
			fmt.Println(errorMessage)
		}

		logger.Fatal().Err(err).Msg("An error occurred when checking if user already exists")
	}

	if existingUserByEmail != nil {
		err := errors.New("A user with the provided email address already exists")

		if envType == "production" {
			fmt.Print(err.Error())
		}

		logger.Fatal().Err(err).Msg("An error occurred when checking if the user already exists")
	}

	userRoleNames := slices.Clone(splitRoles)
	user := &models.User{
		Username:  validatorData.Username,
		Email:     validatorData.Email,
		FirstName: validatorData.FirstName,
		LastName:  validatorData.LastName,
	}

	if *superuser || userCount == 0 {
		userRoleNames = append(userRoleNames, "superuser")
	} else if *admin {
		userRoleNames = append(userRoleNames, "admin")
	} else if len(splitRoles) == 0 {
		userRoleNames = append(userRoleNames, "user")
	}

	if err := user.HashPassword(validatorData.Password); err != nil {
		if envType == "production" {
			fmt.Println(errorMessage)
		}

		logger.Fatal().Err(err).Msg("An error occurred while hashing user password")
	}

	userRoles := make([]*models.Role, 0)
	if err := dbInstance.Select("id").Where("name IN ?", userRoleNames).Find(&userRoles).Error; err != nil {
		if envType == "production" {
			fmt.Println(err)
		}

		logger.Fatal().Err(err).Msg("An error occurred when fetching roles for user creation")
	}

	userRoleIDs := make([]string, len(userRoles))

	for index := range userRoles {
		userRoleIDs[index] = userRoles[index].ID
	}

	err2 := userRepo.CreateWithRoles(user, userRoleIDs, repository.QueryOptions{Omit: []string{"Sessions"}}, repository.QueryOptions{Omit: []string{"Permissions"}})
	if err2 != nil {
		if envType == "production" {
			fmt.Println(errorMessage)
		}

		logger.Fatal().Err(err).Msg("An error occurred during the user creation process")
	}

	fmt.Printf("user %s was successfully created!!!\n", user.Username)
}
