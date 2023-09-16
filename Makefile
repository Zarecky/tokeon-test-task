-include .env

MIGRATIONS_DIR   = ./sql/migrations/
SHELL            := /bin/sh
GOBIN            ?= $(GOPATH)/bin
PATH             := $(GOBIN):$(PATH)
GO               = go
TARGET_DIR       ?= $(PWD)/.build
POSTGRES_DSN_SSL	 = postgres://$(POSTGRES_USER):$(POSTGRES_PASS)@$(POSTGRES_ADDR)/$(POSTGRES_DB_NAME)
POSTGRES_DSN	 = $(POSTGRES_DSN_SSL)?sslmode=disable

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

gensql:
	pgxgen crud -c "$(POSTGRES_DSN)"
	pgxgen sqlc generate

migrate:
	migrate -path "$(MIGRATIONS_DIR)" -database "$(POSTGRES_DSN_SSL)" $(filter-out $@,$(MAKECMDGOALS))

migrate-local:
	migrate -path "$(MIGRATIONS_DIR)" -database "$(POSTGRES_DSN)" $(filter-out $@,$(MAKECMDGOALS))

db-create-migration:
	migrate create -ext sql -dir "$(MIGRATIONS_DIR)" $(filter-out $@,$(MAKECMDGOALS))

gensql:
	pgxgen crud
	pgxgen sqlc generate

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

