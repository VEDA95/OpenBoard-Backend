package main

import (
	"VEDA95/open_board/api/internal/config"
	"VEDA95/open_board/api/internal/db"
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/routes"
	"VEDA95/open_board/api/internal/http/validators"
	applogger "VEDA95/open_board/api/internal/log"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"os"
)

func main() {
	if err := config.LoadEnvConfigs("./env"); err != nil {
		log.Fatal(err)
	}

	if err := applogger.InitializeLogger(); err != nil {
		log.Fatal(err)
	}

	if err := db.InitializeDBInstance(); err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := db.Instance.Close()

		if err != nil {
			applogger.Logger.Fatal().Err(err).Msg("failed to close database instance")
		}
	}()

	validators.InitializeValidatorInstance()

	port := os.Getenv("PORT")

	if len(port) == 0 {
		applogger.Logger.Fatal().Msg("environment variable PORT is not set")
	}

	host := os.Getenv("HOST")
	var hostString string

	if len(host) == 0 {
		hostString = fmt.Sprintf("127.0.0.1:%s", port)

	} else if host == "0.0.0.0" {
		hostString = fmt.Sprintf(":%s", port)

	} else {
		hostString = fmt.Sprintf("%s:%s", host, port)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: errors.ErrorHandler,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})
	apiGroup := app.Group("/api")

	app.Use(fiberzerolog.New(fiberzerolog.Config{Logger: &applogger.Logger}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PATCH, DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	apiGroup.Get("/", routes.IndexGET)

	if err := app.Listen(hostString); err != nil {
		applogger.Logger.Fatal().Err(err).Msg("Error occurred while running the server")
	}
}
