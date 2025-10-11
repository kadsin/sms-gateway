# ------- APP

init:
	@echo - Installing dependencies ...
	go mod download

test:
	@go test -p 1 \
		$(if $(path), $(path), $(if $(scope), ./tests/$(scope)/*, ./...)) \
		$(if $(race), -race) \
		$(if $(verbose), -v) \
		-count $(if $(count), $(count), 1) \
		-run ^.*$(filter).* \
		| grep -a --line-buffered -v -e "no test files" -e "no tests to run"

# ------- DOCKER

docker\:up:
	cd deployments && docker compose --env-file ../.env up -d --build --remove-orphans

docker\:down:
	cd deployments && docker compose down

# ---------------------- Migration tools -------------------------

-include .env

GOOSE_NAME:= database-migrator-cli
BUILD_GOOSE:= go build -o ./build/$(GOOSE_NAME) ./cmd/$(GOOSE_NAME)
GOOSE_CMD := $(BUILD_GOOSE) && ./build/$(GOOSE_NAME)

migrate\:create:
	@$(GOOSE_CMD) create "$(name)" go

migrate:
	@$(GOOSE_CMD) up

migrate\:rollback:
	@$(GOOSE_CMD) down

migrate\:version:
	@$(GOOSE_CMD) version

migrate\:status:
	@$(GOOSE_CMD) status
