PROJECTNAME="Bird Bot"
PROJECT_BIN="birdbot"
VERSION="DEV"
BUILD_NUMBER:=$(shell git rev-parse --short HEAD)

# Go related variables.
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/build
GOFILES=$(wildcard *.go)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

go-full-build: go-clean go-get go-build

go-build:
	@echo "  >  Building binary..."
	@mkdir -p $(GOBIN)
	@GOOS=linux CGO_ENABLED=1 go build -ldflags "-X github.com/yeslayla/birdbot/app.Version=$(VERSION) -X github.com/yeslayla/birdbot/app.Build=$(BUILD_NUMBER)" -o $(GOBIN)/$(PROJECT_BIN) $(GOFILES)
	@chmod 755 $(GOBIN)/$(PROJECT_BIN)

go-generate:
	@echo "  >  Generating dependency files..."
	@go generate $(generate)

go-get:
	@go env -w GOPRIVATE=github.com/meteoritesolutions
	@echo "  >  Checking if there is any missing dependencies..."
	@go get $(get)

go-install:
	@echo "  >  Running go install..."
	@go install $(GOFILES)

go-clean:
	@echo "  >  Cleaning build cache"
	@go clean

go-test: clean
	@echo "  >  Running tests..."
	@go test -coverprofile=coverage.out ./*/

go-run:
	@echo "  >  Running ${PROJECTNAME}"
	@-(cd $(GOBIN); ./$(PROJECT_BIN))

docker-build:
	@docker build . -t yeslayla/birdbot:latest

docker-run: docker-build
	@docker run -it -v `pwd`/build:/etc/birdbot yeslayla/birdbot:latest

docker-push: docker-build
	@docker push yeslayla/birdbot:latest

## install: Download and install dependencies
install: go-get

# clean: Runs go clean
clean: go-clean

## full-build: cleans project, installs dependencies, and builds project
full-build: go-full-build

## build: Runs go build
build: go-build

## package: Builds lambda zip
package: go-full-build
	@echo "  >  Zipping package..."
	@cd $(GOBIN) && zip $(PROJECTNAME).zip $(PROJECTNAME)

## clean: Runs go clean
clean:
	@rm -rf build

## run: full-builds and executes project binary
run: go-build go-run

## test: Run unit tests
test: go-test

## help: Displays help text for make commands
.DEFAULT_GOAL := help
all: help
help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'