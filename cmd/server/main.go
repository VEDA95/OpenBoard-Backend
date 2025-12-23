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

	generalSettingsRepo := repository.NewGeneralSettingsRepository(dbInstance)
	authSettingsRepo := repository.NewAuthSettingsRepository(dbInstance)
	emailSettingsRepo := repository.NewEmailSettingsRepository(dbInstance)
	emailService, err := service.NewEmailService(emailRepo, emailSettingsRepo)
	if err != nil {
		errMessage := err.Error()
		if !strings.Contains(errMessage, "required email settings are missing") && !strings.Contains(errMessage, "valid email address for the sender is required") {
			applogger.Global.Fatal().Err(err).Msg("")
		}

		applogger.Global.Warn().Err(err).Msg("")
	}

	userRepo := repository.NewUserRepository(dbInstance)
	sessionRepo := repository.NewSessionRepository(dbInstance)
	roleRepo := repository.NewRoleRepository(dbInstance)
	permissionRepo := repository.NewPermissionRepository(dbInstance)
	passwordResetRepo := repository.NewPasswordResetRepository(dbInstance)
	workspaceRepo := repository.NewWorkspaceRepository(dbInstance)
	boardRepo := repository.NewBoardRepository(dbInstance)
	listRepo := repository.NewListRepository(dbInstance)
	cardRepo := repository.NewCardRepository(dbInstance)
	checkListItemRepo := repository.NewCheckListItemRepository(dbInstance)
	labelRepo := repository.NewLabelRepository(dbInstance)
	commentRepo := repository.NewCommentRepository(dbInstance)
	activityRepo := repository.NewActivityRepository(dbInstance)
	fileUploadRepo := repository.NewFileUploadRepository(dbInstance)
	settingsService := service.NewSettingsService(generalSettingsRepo, authSettingsRepo, emailSettingsRepo)
	userService := service.NewUserService(userRepo, roleRepo)
	authService := service.NewAuthService(sessionRepo, userRepo, passwordResetRepo, roleRepo, authSettingsRepo)
	roleService := service.NewRoleService(roleRepo)
	permissionService := service.NewPermissionService(permissionRepo)
	workspaceService := service.NewWorkspaceService(workspaceRepo)
	boardService := service.NewBoardService(boardRepo)
	listService := service.NewListService(listRepo)
	cardService := service.NewCardService(cardRepo)
	labelService := service.NewLabelService(labelRepo, boardRepo)
	checkListItemService := service.NewCheckListItemService(checkListItemRepo, cardRepo, activityRepo)
	commentService := service.NewCommentService(commentRepo, cardRepo)
	fileUploadService := service.NewFileUploadService(fileUploadRepo, userRepo, cardRepo)
	validator := validators.NewValidator()
	settingsHandler := routes.NewSettingsHandler(settingsService, emailService, validator)
	userHandler := routes.NewUserHandler(userService, fileUploadService, validator)
	authHandler := routes.NewAuthHandler(authService, userService, emailService, validator)
	roleHandler := routes.NewRoleHandler(roleService, validator)
	permissionHandler := routes.NewPermissionHandler(permissionService, validator)
	workspaceHandler := routes.NewWorkspaceHandler(workspaceService, validator)
	boardHandler := routes.NewBoardHandler(boardService, validator)
	listHandler := routes.NewListHandler(listService, validator)
	cardHandler := routes.NewCardHandler(cardService, validator)
	labelHandler := routes.NewLabelHandler(labelService, validator)
	checkListItemHandler := routes.NewCheckListItemHandler(checkListItemService, validator)
	commentHandler := routes.NewCommentHandler(commentService, validator)
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
	authGroup.Post("/register", authHandler.Register)
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
	authGroup.Get("/roles", authMiddleware.RequireAuthentication(), roleHandler.GET)
	authGroup.Post("/roles", authMiddleware.RequireAuthentication(), roleHandler.POST)
	authGroup.Get("/roles/:id", authMiddleware.RequireAuthentication(), roleHandler.GETByID)
	authGroup.Patch("/roles/:id", authMiddleware.RequireAuthentication(), roleHandler.PATCH)
	authGroup.Delete("/roles/:id", authMiddleware.RequireAuthentication(), roleHandler.DELETE)
	authGroup.Get("/permissions", authMiddleware.RequireAuthentication(), permissionHandler.GET)
	authGroup.Post("/permissions", authMiddleware.RequireAuthentication(), permissionHandler.POST)
	authGroup.Get("/permissions/:id", authMiddleware.RequireAuthentication(), permissionHandler.GETByID)
	authGroup.Patch("/permissions/:id", authMiddleware.RequireAuthentication(), permissionHandler.PATCH)
	authGroup.Delete("/permissions/:id", authMiddleware.RequireAuthentication(), permissionHandler.DELETE)
	apiGroup.Get("/settings/:type", authMiddleware.RequireAuthentication(), settingsHandler.GET)
	apiGroup.Get("/settings/:type/public", settingsHandler.GETPublic)
	apiGroup.Patch(
		"/settings/:type",
		authMiddleware.RequireAuthentication(),
		authMiddleware.RequireAuthorization("system:settings"),
		settingsHandler.PATCH,
	)
	apiGroup.Get("/users", authMiddleware.RequireAuthentication(), userHandler.GET)
	apiGroup.Post("/users", authMiddleware.RequireAuthentication(), userHandler.POST)
	apiGroup.Get("/users/:id", authMiddleware.RequireAuthentication(), userHandler.GETByID)
	apiGroup.Patch("/users/:id", authMiddleware.RequireAuthentication(), userHandler.PATCH)
	apiGroup.Delete("/users/:id", authMiddleware.RequireAuthentication(), userHandler.DELETE)
	apiGroup.Get("/workspaces", authMiddleware.RequireAuthentication(), workspaceHandler.GET)
	apiGroup.Post("/workspaces", authMiddleware.RequireAuthentication(), workspaceHandler.POST)
	apiGroup.Get("/workspaces/:id", authMiddleware.RequireAuthentication(), workspaceHandler.GETByID)
	apiGroup.Patch("/workspaces/:id", authMiddleware.RequireAuthentication(), workspaceHandler.PATCH)
	apiGroup.Delete("/workspaces/:id", authMiddleware.RequireAuthentication(), workspaceHandler.DELETE)
	apiGroup.Get("/boards", authMiddleware.RequireAuthentication(), boardHandler.GET)
	apiGroup.Post("/boards", authMiddleware.RequireAuthentication(), boardHandler.POST)
	apiGroup.Get("/boards/:id", authMiddleware.RequireAuthentication(), boardHandler.GETByID)
	apiGroup.Patch("/boards/:id", authMiddleware.RequireAuthentication(), boardHandler.PATCH)
	apiGroup.Delete("/boards/:id", authMiddleware.RequireAuthentication(), boardHandler.DELETE)
	apiGroup.Get("/lists", authMiddleware.RequireAuthentication(), listHandler.GET)
	apiGroup.Post("/lists", authMiddleware.RequireAuthentication(), listHandler.POST)
	apiGroup.Get("/lists/:id", authMiddleware.RequireAuthentication(), listHandler.GETByID)
	apiGroup.Patch("/lists/:id", authMiddleware.RequireAuthentication(), listHandler.PATCH)
	apiGroup.Delete("/lists/:id", authMiddleware.RequireAuthentication(), listHandler.DELETE)
	apiGroup.Get("/cards", authMiddleware.RequireAuthentication(), cardHandler.GET)
	apiGroup.Post("/cards", authMiddleware.RequireAuthentication(), cardHandler.POST)
	apiGroup.Get("/cards/:id", authMiddleware.RequireAuthentication(), cardHandler.GETByID)
	apiGroup.Patch("/cards/:id", authMiddleware.RequireAuthentication(), cardHandler.PATCH)
	apiGroup.Patch("/cards/:id/move", authMiddleware.RequireAuthentication(), cardHandler.Move)
	apiGroup.Delete("/cards/:id", authMiddleware.RequireAuthentication(), cardHandler.DELETE)
	apiGroup.Get("/labels", authMiddleware.RequireAuthentication(), labelHandler.GET)
	apiGroup.Post("/labels", authMiddleware.RequireAuthentication(), labelHandler.POST)
	apiGroup.Get("/labels/:id", authMiddleware.RequireAuthentication(), labelHandler.GETByID)
	apiGroup.Patch("/labels/:id", authMiddleware.RequireAuthentication(), labelHandler.PATCH)
	apiGroup.Delete("/labels/:id", authMiddleware.RequireAuthentication(), labelHandler.DELETE)
	apiGroup.Get("/check_list_items", authMiddleware.RequireAuthentication(), checkListItemHandler.GET)
	apiGroup.Post("/check_list_items", authMiddleware.RequireAuthentication(), checkListItemHandler.POST)
	apiGroup.Get("/check_list_items/:id", authMiddleware.RequireAuthentication(), checkListItemHandler.GETByID)
	apiGroup.Patch("/check_list_items/:id", authMiddleware.RequireAuthentication(), checkListItemHandler.PATCH)
	apiGroup.Delete("/check_list_items/:id", authMiddleware.RequireAuthentication(), checkListItemHandler.DELETE)
	apiGroup.Get("/comments", authMiddleware.RequireAuthentication(), commentHandler.GET)
	apiGroup.Post("/comments", authMiddleware.RequireAuthentication(), commentHandler.POST)
	apiGroup.Get("/comments/:id", authMiddleware.RequireAuthentication(), commentHandler.GETByID)
	apiGroup.Patch("/comments/:id", authMiddleware.RequireAuthentication(), commentHandler.PATCH)
	apiGroup.Delete("/comments/:id", authMiddleware.RequireAuthentication(), commentHandler.DELETE)

	if err := app.Listen(hostString); err != nil {
		applogger.Global.Fatal().Err(err).Msg("Error occurred while running the server")
	}
}
