-include .env

SHELL            := /bin/sh
GOBIN            ?= $(GOPATH)/bin
PATH             := $(GOBIN):$(PATH)
GO               = go
TARGET_DIR       ?= $(PWD)/.build

ifeq ($(DELVE_ENABLED),true)
GCFLAGS	= -gcflags 'all=-N -l'
endif

# Setting path to installed go tools if $GOPATH is empty
ifeq ($(GOPATH),)
GOBIN ?=$(GOBIN)
endif


.PHONY: start
start:
	go run ./cmd/app/main.go


.PHONY: run
run:
	$(TARGET_DIR)/cmd/app

.PHONY: build
build:
	$(info $(M) building application...)
	@GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build $(GCFLAGS) $(LDFLAGS) -o $(TARGET_DIR)/cmd/app ./cmd/app/*.go

.PHONY: watch
watch: ## Run binaries that rebuild themselves on changes
	$(info $(M) run...)
	CONSUL_STAND_NAME=local air -c $(PWD)/.air.conf

genswagger:
	swag init -g /cmd/app/main.go --parseDependency --parseInternal

compose:
	docker compose up -d


install-tools: # Install tools needed for development
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/cosmtrek/air@latest
	go install github.com/oligot/go-mod-upgrade@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/tkcrm/pgxgen/cmd/pgxgen@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/fullstorydev/grpcui/cmd/grpcui@latest


%:
	@:

