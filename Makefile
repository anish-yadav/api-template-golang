VERSION := local
CLI_OUTPUT := lms.exe
SERVER_OUTPUT := server.exe
BIN_DIR := bin

download:
	@echo "go mod download"
	@go mod download

build-win :
	@echo "Building for windows"
	@echo "\033[0;31m @GOOS=windows GOARCH=amd64 go build -o ${BIN_DIR}/${SERVER_OUTPUT} cmd/server/main.go"
	@GOOS=windows GOARCH=amd64 go build -o ${BIN_DIR}/${SERVER_OUTPUT} cmd/server/main.go
	@echo "\033[1;33m @GOOS=windows GOARCH=amd64 go build -o ${BIN_DIR}/${CLI_OUTPUT} cmd/root/main.go"
	@GOOS=windows GOARCH=amd64 go build -o ${BIN_DIR}/${CLI_OUTPUT} cmd/root/main.go

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo