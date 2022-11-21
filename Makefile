BIN_NAME=worklog

test:
	@go clean -testcache ./...
	@make --no-print-directory test-unit
	@make --no-print-directory test-e2e

test-repeat:
	@go clean -testcache ./...
	@make --no-print-directory test-unit ARGS="-count 50"
	@make --no-print-directory test-e2e

test-unit:
	@echo "Running unit tests"
	go test $(ARGS) ./cli ./helpers ./model ./repository ./service
	@echo "Unit tests passed"

test-e2e:
	@echo "Running end to end tests"
	@make --no-print-directory build DIR="./e2e"

	# Run once so all files are setup
	@./e2e/worklog configure > /tmp/dump
	@rm -f /tmp/e2e.db

	commander test ./e2e/root.test.yaml
	commander test ./e2e/configure.test.yaml
	commander test ./e2e/create.test.yaml
	commander test ./e2e/print.test.yaml
	commander test ./e2e/edit.test.yaml
	@make --no-print-directory clean
	@echo "e2e tests passed"

format:
	@echo "Running format checks"
	go fmt ./...
	go vet ./...
	@echo "Format checks passed"

build:
	@echo "Building to $(DIR)"
	@go build -ldflags="-w -s" -o "$(DIR)/$(BIN_NAME)" exec/cli/main.go
	@go build -ldflags="-w -s" -o "$(DIR)/$(BIN_NAME)-server" exec/server/main.go
	@echo "Built to $(DIR)"

build-local:
	@make --no-print-directory build DIR="/tmp"

compile:
	@echo "Compiling for OS and Platform"
	goreleaser release --snapshot --rm-dist
	@echo "Built all binaries"

deps:
	@echo "Checking and downloading dependencies"
	go mod tidy
	go mod download
	go mod verify
	@echo "Passed all dependency checks"

clean:
	@echo "Cleaning"
	@rm -rf dist/ ./e2e/$(BIN_NAME) ./e2e/$(BIN_NAME)-server
