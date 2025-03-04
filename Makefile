GO_BUILD = go build -v -o
BUILD_DIR = ./build
API_ENTRYPOINT = ./cmd/server/main.go
MIGRATION_ENTRYPOINT = ./cmd/migrate/main.go

install: build_migration_app build_api

build_api:
	$(GO_BUILD) $(BUILD_DIR)/backend_api $(API_ENTRYPOINT)

build_migration_app:
	$(GO_BUILD) $(BUILD_DIR)/migration_app $(MIGRATION_ENTRYPOINT)

run_dev_api:
	air -c ./deployment/dev/.air.linux.conf

run_dev_migration_app:
	go run $(MIGRATION_ENTRYPOINT) $(ARGS)

clean:
	go clean -cache -modcache
	rm -rf ./build/*%                                                                                                                                                                                                              veda@Stefans-MacBook-Pro api %
