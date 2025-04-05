envfile := .env
-include $(envfile)
export $(shell sed 's/=.*//' $(envfile))

# HELP
.PHONY: help

help: ## Usage: make <option>
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

init: ## APP init configuration.
	if [ ! -f .env ]; then cp -n .env.dist .env; echo "CONTAINER_EXEC_USER_ID=`id -u`" >> .env; echo "CONTAINER_USERNAME=${USER}" >> .env; fi
	#docker compose build xm_app_${APP_ENV};

bash: ## APP Bash.
	docker compose exec xm_app_${APP_ENV} bash

trace-vars: ## APP Trace vars.
	docker compose exec xm_app_${APP_ENV} expvarmon -ports="${XM_DEBUG_PORT}"

up: ## APP Start.
	docker compose up -d --build --remove-orphans xm_app_${APP_ENV};

down: ## APP Stop.
	docker compose down --remove-orphans;

ps: ## APP Processes.
	docker compose ps;

logs: ## APP Logs.
	docker compose logs -f;

example: ## APP Example.
	docker compose exec xm_app_${APP_ENV} /srv/xm/bin/app api-example

lint: ## APP Lint.
	docker run -t --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v2.0.2 golangci-lint run

test: down ## APP Test - and show coverage
	docker compose up -d --build --remove-orphans xm_app_dev;
	docker compose exec xm_app_dev go test -race -cpu 24 -cover -coverprofile=data/test/coverage.out ./internal/...;
	docker compose run xm_app_dev go tool cover -func=data/test/coverage.out;
	make down;


# not integrated with docker
trace-pprof-allocs:
	go tool pprof -http :8090 http://localhost:4000/debug/pprof/allocs?debug=1
