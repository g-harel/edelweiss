.PHONY: help check install run setup test check-go check-dep check-kubectl check-minikube

GO ?= $(shell command -v go 2> /dev/null)
DEP ?= $(shell command -v dep 2> /dev/null)
KUBECTL ?= $(shell command -v kubectl 2> /dev/null)
MINIKUBE ?= $(shell command -v minikube 2> /dev/null)

help:
	@echo "Make Commands:"
	@echo "    help - list commands"
	@echo "    deps - check if required dependencies are installed"

check: check-go check-dep check-kubectl check-minikube
	@echo checking dependencies ...

install: check-dep
	@$(DEP) ensure

run: check-go
	@docker-compose up -d
	@$(GO) run ./services/gateway/main.go

setup: check-minikube
	@$(MINIKUBE) start
	@$(MINIKUBE) addons enable registry

test: check-go
	@$(GO) test -race ./...

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
