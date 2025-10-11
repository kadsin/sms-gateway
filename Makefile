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
