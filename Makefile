.DEFAULT_GOAL := help
RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)

start=reflex -r '(\.go$|go\.mod)' -R .idea/ -s -d none $(2) -- sh -c 'make build && $(or $(value 1), /usr/bin/musicd serve)'

define setup_env
	$(eval ENV_FILE := .env$(1))
	$(eval include .env$(1))
	$(eval export sed 's/=.*//' .env$(1))
	$(eval export PGUSER=${DB_USER})
	$(eval export PGPASSWORD=${DB_PASSWORD})
	$(eval export PGHOST=${DB_HOST})
	$(eval export PGPORT=${DB_PORT})
	$(eval export PGDATABASE=${DB_NAME})
endef

define generate_migration_name
	$(if $(1),$(eval MIGRATION_NAME := $(shell dirname $1)/$(shell date -u +'%Y%m%d%H%M%S')_$(shell basename $1)),)
endef

####################################################################################################
## MAIN COMMANDS
####################################################################################################
help: ## Commands list
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

build:
	go build -buildvcs=false -o /usr/bin/musicd cmd/main.go

start:
	$(call start)

generate: ## Generate code 
	go generate ./...

analyze: ## Run static analyzer
	test -s ./bin/golangci-lint || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.52.2
	./bin/golangci-lint run -c ./.golangci.yaml ./...

test: ## Run tests
	go test ./internal/...

.PHONY: sqitch-link 
sqitch-link:
	rm -rf sqitch/deploy sqitch/verify sqitch/revert
	mkdir -p sqitch/deploy sqitch/verify sqitch/revert
	find internal/ -type d -name storage | while read module; do \
		module_name=$$(dirname $$module | xargs basename); \
		mkdir -p $$module/sqitch/deploy ; \
		ln -sf ../../$$module/sqitch/deploy ./sqitch/deploy/$$module_name; \
		mkdir -p $$module/sqitch/verify ; \
		ln -sf  ../../$$module/sqitch/verify ./sqitch/verify/$$module_name; \
		mkdir -p $$module/sqitch/revert ; \
		ln -sf ../../$$module/sqitch/revert ./sqitch/revert/$$module_name; \
	done

.PHONY: sqitch
sqitch: sqitch-link ## Run sqitch
	$(call setup_env,)
	TZ=UTC sqitch $(RUN_ARGS)

.PHONY: db-migrate
db-migrate: ## Run migrations 
	$(call setup_env,)
	TZ=UTC sqitch deploy $(RUN_ARGS)
	$(MAKE) db-generate 

.PHONY: check-migration
check-migration: ## Run migrations on test environment, then rollback and migrate again
	$(MAKE) db-migrate-test
	$(MAKE) db-rollback-test @HEAD^
	$(MAKE) db-migrate-test

db-add: sqitch-link  ## Add a new migration
	$(call setup_env,)
	$(call generate_migration_name,$(word 1,$(RUN_ARGS)))
	TZ=UTC sqitch add $(MIGRATION_NAME)

db-rollback: sqitch-link ## Rollback database migrations over the real DB
	$(call setup_env,)
	TZ=UTC sqitch revert $(RUN_ARGS)

db-generate: ## Generate DTO and DAO for modules
	find . -path './internal/*/sqlc.yaml' -exec sqlc -f '{}' generate ';'

