# ------- APP

init:
	@echo - Installing dependencies ...
	go mod download

server:
	@go build -o ./build/server ./cmd/server
	@./build/server

test\:help:
	@printf "\
	\n\
	To run all tests:\
	\n\
		make test\
	\n\n\
	To test a path:\
	\n\
		make test path='path/to/your_test(s)'\
	\n\n\
	To test a scope in tests directory:\
	\n\
		make test scope=server\
	\n\n\
	To test in race mode:\
	\n\
		make test race=t\
	\n\n\
	To run tests by a filter:\
	\n\
		make test filter='UserBalance'\
	\n\n\
	To run tests x times:\
	\n\
		make test count=10\
	\n\n\
	To run tests with detail:\
	\n\
		make test verbose=t\
	"

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

DB_GOOSE_NAME:= database-migrator-cli
BUILD_DB_GOOSE:= go build -o ./build/$(DB_GOOSE_NAME) ./cmd/$(DB_GOOSE_NAME)
DB_GOOSE_CMD := $(BUILD_DB_GOOSE) && ./build/$(DB_GOOSE_NAME)

migrate\:db\:create:
	@$(DB_GOOSE_CMD) create "$(name)" go

migrate\:db:
	@$(DB_GOOSE_CMD) up

migrate\:db\:rollback:
	@$(DB_GOOSE_CMD) down

migrate\:db\:version:
	@$(DB_GOOSE_CMD) version

migrate\:db\:status:
	@$(DB_GOOSE_CMD) status


ANALYTICS_GOOSE_NAME:= analytics-migrator-cli
BUILD_ANALYTICS_GOOSE:= go build -o ./build/$(ANALYTICS_GOOSE_NAME) ./cmd/$(ANALYTICS_GOOSE_NAME)
ANALYTICS_GOOSE_CMD := $(BUILD_ANALYTICS_GOOSE) && ./build/$(ANALYTICS_GOOSE_NAME)

migrate\:analytics\:create:
	@$(ANALYTICS_GOOSE_CMD) create "$(name)" go

migrate\:analytics:
	@$(ANALYTICS_GOOSE_CMD) up

migrate\:analytics\:rollback:
	@$(ANALYTICS_GOOSE_CMD) down

migrate\:analytics\:version:
	@$(ANALYTICS_GOOSE_CMD) version

migrate\:analytics\:status:
	@$(ANALYTICS_GOOSE_CMD) status
