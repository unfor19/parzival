
help:                ## Available make commands
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's~:~~' | sed -e 's~##~~'

usage: help         

build:
	@go build

up-localstack:       ## Run localstack in Docker Compose
	@docker-compose -p parzival -f docker-compose-localstack.yml up --detach

down-localstack:     ## Stop localstack in Docker Compose
	@docker-compose -p parzival -f docker-compose-localstack.yml down

clean-localstack:    ## Clean localstack in Docker Compose
	@docker-compose -p parzival -f docker-compose-localstack.yml down -v --remove-orphans
	@docker rm -f localstack 2>/dev/null || true
	@rm -rf .localstack 2>/dev/null || true

test:                ## Run tests
	@./scripts/tests.sh
