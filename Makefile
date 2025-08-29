GO_BUILD = go build -v -o
GO_RUN = go run
BUILD_DIR = ./build
API_ENTRYPOINT = ./cmd/server/main.go
MIGRATION_ENTRYPOINT = ./cmd/migrate/main.go
SEED_ENTRYPOINT = ./cmd/admin/seed/main.go
CREATE_USER_ENTRYPOINT = ./cmd/admin/createuser/main.go
RESET_PASSWORD_ENTRYPOINT = ./cmd/admin/resetpassword/main.go

install: build_seed_program build_user_creation_program build_password_reset_program build_swagger_documentation build_api

build_api:
	$(GO_BUILD) $(BUILD_DIR)/backend_api $(API_ENTRYPOINT)

build_seed_program:
	$(GO_BUILD) $(BUILD_DIR)/seed $(SEED_ENTRYPOINT)

build_user_creation_program:
	$(GO_BUILD) $(BUILD_DIR)/create_user $(CREATE_USER_ENTRYPOINT)

build_password_reset_program:
	$(GO_BUILD) $(BUILD_DIR)/reset_password $(RESET_PASSWORD_ENTRYPOINT)

run_dev_api:
	air -c ./deployment/dev/.air.linux.conf

run_dev_seed_program:
	$(GO_RUN) $(SEED_ENTRYPOINT)

run_dev_create_user_program:
	$(GO_RUN) $(CREATE_USER_ENTRYPOINT)

run_dev_reset_password_program:
	$(GO_RUN) $(RESET_PASSWORD_ENTRYPOINT)

build_swagger_documentation:
	swag init --parseDependency -g $(API_ENTRYPOINT)

clean:
	go clean -cache -modcache
	rm -rf ./build/*%                                                                                                                                                                                                              veda@Stefans-MacBook-Pro api %
