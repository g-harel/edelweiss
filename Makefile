.PHONY: help check install run setup test check-go check-dep check-kubectl check-minikube gateway

GO ?= $(shell command -v go 2> /dev/null)
DEP ?= $(shell command -v dep 2> /dev/null)
KUBECTL ?= $(shell command -v kubectl 2> /dev/null)
MINIKUBE ?= $(shell command -v minikube 2> /dev/null)

help:
	@echo "Make Commands:"
	@echo "    help - list commands"
	@echo "    deps - check if required dependencies are installed"
	@echo "    test - run tests"
	@echo "    run  - starts all services"

check: check-go check-dep check-kubectl check-minikube
	@echo checking dependencies ...

install: check-dep
	@$(DEP) ensure

run: check-go
	@docker-compose up -d
	@$(GO) run ./services/gateway/main.go

test: check-go
	@$(GO) test -race ./...

gateway:
	@GOOS=linux GOARCH=amd64 go build -o=build/gateway ./services/gateway/main.go

# checks for required commands

check-go:
ifndef GO
	$(error "could not locate go")
endif

check-dep:
ifndef DEP
	$(error "could not locate dep")
endif

check-kubectl:
ifndef KUBECTL
	$(error "could not locate kubectl")
endif

check-minikube:
ifndef MINIKUBE
	$(error "could not locate minikube")
endif
