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
	@go test $(ARGS) ./cmd ./helpers ./model ./repository ./service
	@echo "Unit tests passed"

test-e2e:
	@echo "Running e2e tests"
	@rm ./e2e/$(BIN_NAME) || true
	@rm ./e2e/e2e.* || true
	@make --no-print-directory build DIR="./e2e"
	@go test ./e2e
	@rm ./e2e/$(BIN_NAME) || true
	@rm ./e2e/e2e.* || true
	@echo "e2e tests passed"

format:
	@echo "Running format checks"
	go fmt ./...
	go vet ./...
	@echo "Format checks passed"

build:
	@go build -ldflags="-w -s" -o "$(DIR)/$(BIN_NAME)" main.go

build-local:
	@make --no-print-directory build DIR="/tmp"

compile:
	@echo "Compiling for OS and Platform"
	@echo "32 bit systems"
	@GOOS=freebsd GOARCH=386 go build -ldflags='-w -s' -o bin/32bit/$(BIN_NAME)-freebsd main.go
	@GOOS=linux GOARCH=386 go build -ldflags='-w -s' -o bin/32bit/$(BIN_NAME)-linux main.go
	@GOOS=windows GOARCH=386 go build -ldflags='-w -s' -o bin/32bit/$(BIN_NAME)-windows.exe main.go

	@echo "64 bit systems"
	@GOOS=freebsd GOARCH=amd64 go build -ldflags='-w -s' -o bin/64bit/$(BIN_NAME)-freebsd main.go
	@GOOS=darwin GOARCH=amd64 go build -ldflags='-w -s' -o bin/64bit/$(BIN_NAME)-darwin main.go
	@GOOS=linux GOARCH=amd64 go build -ldflags='-w -s' -o bin/64bit/$(BIN_NAME)-linux main.go
	@GOOS=windows GOARCH=amd64 go build -ldflags='-w -s' -o bin/64bit/$(BIN_NAME)-windows.exe main.go
	@echo "Built all binaries"

deps:
	@echo "Checking and downloading dependencies"
	go mod tidy
	go mod download
	go mod verify
	@echo "Passed all dependency checks"
