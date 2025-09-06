package main

import (
	_ "VEDA95/open_board/api/docs"
	"VEDA95/open_board/api/internal/config"
	"VEDA95/open_board/api/internal/db"
	"VEDA95/open_board/api/internal/db/repository"
	"VEDA95/open_board/api/internal/errors"
	"VEDA95/open_board/api/internal/http/middleware"
	"VEDA95/open_board/api/internal/http/routes"
	"VEDA95/open_board/api/internal/http/validators"
	applogger "VEDA95/open_board/api/internal/log"
	"VEDA95/open_board/api/internal/service"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
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

	port := os.Getenv("PORT")
	if len(port) == 0 {
		log.Fatal("environment variable PORT is not set")
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

	if err := applogger.InitializeGlobal(); err != nil {
		log.Fatal(err)
	}

	dbInstance, err := db.NewDB()
	if err != nil {
		applogger.Global.Fatal().Err(err).Msg("")
	}

	emailRepo, err := repository.NewEmailRepository()
	if err != nil {
		applogger.Global.Fatal().Err(err).Msg("")
	}

	emailService, err := service.NewEmailService(emailRepo)
	if err != nil {
		if !strings.Contains(err.Error(), "required email settings are missing") {
			applogger.Global.Fatal().Err(err).Msg("")
		}

		applogger.Global.Warn().Err(err).Msg("")
	}

	userRepo := repository.NewUserRepository(dbInstance)
	sessionRepo := repository.NewSessionRepository(dbInstance)
	roleRepo := repository.NewRoleRepository(dbInstance)
	permissionRepo := repository.NewPermissionRepository(dbInstance)
	passwordResetRepo := repository.NewPasswordResetRepository(dbInstance)
	userService := service.NewUserService(userRepo, roleRepo)
	authService := service.NewAuthService(sessionRepo, userRepo, passwordResetRepo)
	roleService := service.NewRoleService(roleRepo)
	permissionService := service.NewPermissionService(permissionRepo)
	validator := validators.NewValidator()
	userHandler := routes.NewUserHandler(userService, validator)
	authHandler := routes.NewAuthHandler(authService, userService, emailService, validator)
	roleHandler := routes.NewRoleHandler(roleService, validator)
	permissionHandler := routes.NewPermissionHandler(permissionService, validator)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	app := fiber.New(fiber.Config{
		ErrorHandler: errors.ErrorHandler,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})
	apiGroup := app.Group("/api")
	authGroup := app.Group("/auth")

	app.Use(fiberzerolog.New(fiberzerolog.Config{Logger: applogger.Global}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))
	app.Get("/swagger/*", swagger.HandlerDefault)
	authGroup.Post("/login", authHandler.LocalLogin)
	authGroup.Post("/refresh", authHandler.LocalRefresh)
	authGroup.Post("/logout/:id", authMiddleware.RequireAuthentication(), authHandler.LocalLogoutById)
	authGroup.Post("/logout", authMiddleware.RequireAuthentication(), authHandler.LocalLogout)
	authGroup.Post("/password_reset", authMiddleware.RequireAuthentication(), authHandler.LocalAuthenticatedPasswordReset)
	authGroup.Post("/password_reset/token", authMiddleware.RequireAuthentication(), authHandler.LocalAuthenticatedPasswordTokenIssuer)
	authGroup.Post("/forgot_password", authHandler.LocalUnauthenticatedPasswordReset)
	authGroup.Post("/forgot_password/token", authHandler.LocalUnauthenticatedPasswordTokenIssuer)
	authGroup.Get("/@me", authMiddleware.RequireAuthentication(), authHandler.UserInfoGET)
	authGroup.Patch("/@me", authMiddleware.RequireAuthentication(), authHandler.UserInfoPATCH)
	authGroup.Delete("/@me", authMiddleware.RequireAuthentication(), authHandler.UserInfoDELETE)
	authGroup.Get("/@me/sessions", authMiddleware.RequireAuthentication(), authHandler.UserSessionsGET)
	authGroup.Get("/roles", roleHandler.GET)
	authGroup.Post("/roles", roleHandler.POST)
	authGroup.Get("/roles/:id", roleHandler.GETByID)
	authGroup.Patch("/roles/:id", roleHandler.PATCH)
	authGroup.Delete("/roles/:id", roleHandler.DELETE)
	authGroup.Get("/permissions", permissionHandler.GET)
	authGroup.Post("/permissions", permissionHandler.POST)
	authGroup.Get("/permissions/:id", permissionHandler.GETByID)
	authGroup.Patch("/permissions/:id", permissionHandler.PATCH)
	authGroup.Delete("/permissions/:id", permissionHandler.DELETE)
	apiGroup.Get("/users", userHandler.GET)
	apiGroup.Post("/users", userHandler.POST)
	apiGroup.Get("/users/:id", userHandler.GETByID)
	apiGroup.Patch("/users/:id", userHandler.PATCH)
	apiGroup.Delete("/users/:id", userHandler.DELETE)

	if err := app.Listen(hostString); err != nil {
		applogger.Global.Fatal().Err(err).Msg("Error occurred while running the server")
	}
}
