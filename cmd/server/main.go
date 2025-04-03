package main

import (
	_ "VEDA95/open_board/api/docs"
	"VEDA95/open_board/api/internal/config"
	"VEDA95/open_board/api/internal/db"
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/middleware"
	"VEDA95/open_board/api/internal/http/routes"
	"VEDA95/open_board/api/internal/http/validators"
	applogger "VEDA95/open_board/api/internal/log"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"log"
	"os"
)

// @Title Open Board Backend API
// @Version 1.0
// @Description This is the REST API backend for the Open Board project...
// @Accept json
// @Produce json
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

	defer db.Instance.Close()

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
	authGroup := app.Group("/auth")

	app.Use(fiberzerolog.New(fiberzerolog.Config{Logger: &applogger.Logger}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PATCH, DELETE",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	app.Get("/swagger/*", swagger.HandlerDefault)
	authGroup.Post("/login", routes.LocalLogin)
	authGroup.Post("/refresh", routes.LocalRefresh)
	authGroup.Post("/logout/:id", middleware.CheckUserAuthentication, routes.LocalLogoutById)
	authGroup.Post("/logout", middleware.CheckUserAuthentication, routes.LocalLogout)
	authGroup.Get("/@me", middleware.CheckUserAuthentication, routes.UserInfoGET)
	authGroup.Patch("/@me", middleware.CheckUserAuthentication, routes.UserInfoPATCH)
	authGroup.Delete("/@me", middleware.CheckUserAuthentication, routes.UserInfoDELETE)
	authGroup.Get("/@me/sessions", middleware.CheckUserAuthentication, routes.UserSessionsGET)
	authGroup.Get("/roles", routes.RolesGET)
	authGroup.Post("/roles", routes.RolesPOST)
	authGroup.Get("/roles/:id", routes.RoleGET)
	authGroup.Patch("/roles/:id", routes.RolePATCH)
	authGroup.Delete("/roles/:id", routes.RoleDELETE)
	authGroup.Get("/permissions", routes.PermissionsGET)
	authGroup.Post("/permissions", routes.PermissionsPOST)
	authGroup.Get("/permissions/:id", routes.PermissionGET)
	authGroup.Patch("/permissions/:id", routes.PermissionPATCH)
	authGroup.Delete("/permissions/:id", routes.PermissionDELETE)
	apiGroup.Get("/", routes.IndexGET)
	apiGroup.Get("/users", routes.UsersGET)
	apiGroup.Post("/users", routes.UsersPOST)
	apiGroup.Get("/users/:id", routes.UserGET)
	apiGroup.Patch("/users/:id", routes.UserPATCH)
	apiGroup.Delete("/users/:id", routes.UserDELETE)

	if err := app.Listen(hostString); err != nil {
		applogger.Logger.Fatal().Err(err).Msg("Error occurred while running the server")
	}
}
