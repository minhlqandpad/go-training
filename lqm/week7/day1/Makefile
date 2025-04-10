up: ## Start the application dependencies
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

down: ## Stop the application dependencies
	@docker-compose down

run:
	go run cmd/app/main.go

mock:
	mockgen -source=internal/domain/repository/tenant.go -destination=internal/infrastructure/db/mock_db/mock_tenant_repository.go -package=mocks

test:
	go test -v -race -cover ./...
