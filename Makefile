include .env

.DEFAULT_GOAL := help
COMPOSE_COMMAND = docker compose --env-file .env -f build/docker-compose-dev.yml -p dbkit

.PHONY: help config build db acceptance up stop dev githooks
help: ## Show this help
	@echo "\033[36mUsage:\033[0m"
	@echo "make TASK"
	@echo
	@echo "\033[36mTasks:\033[0m"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[33m%-30s\033[0m \r\033[32C%s\n", $$1, $$2}'

# DOCKER TASKS
config: ## Validate and view the Compose file.
	$(COMPOSE_COMMAND) config
dev: ## Build and start dev environment
	$(COMPOSE_COMMAND) up -d --force-recreate --remove-orphans --build db adminer --wait
stop: ## Stops running containers without removing them [s=services]
	$(COMPOSE_COMMAND) stop db
down: ## Stops containers and removes containers, networks, volumes, and images created by up
	$(COMPOSE_COMMAND) down --rmi all --remove-orphans
backup: ## Backup database
	docker exec dbkit-db mysqldump -u root --password=${MYSQL_ROOT_PASSWORD} ${MYSQL_DATABASE} > ./tests/acceptance/data/acceptance.sql
githooks: ## Install git hooks
	chmod +x githooks/* && cp -rf githooks/* .git/hooks/
