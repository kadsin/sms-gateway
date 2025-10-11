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
