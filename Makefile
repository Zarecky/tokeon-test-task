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

%:
	@:

