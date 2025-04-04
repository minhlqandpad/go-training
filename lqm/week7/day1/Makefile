up: ## Start the application dependencies
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

down: ## Stop the application dependencies
	@docker-compose down

run:
	go run cmd/app/main.go

mock:
	mockgen -source=internal/domain/infrastructure/db/tenant.go -destination=internal/domain/infrastructure/db/mock_db/mock_tenant.go -package=mocks TenantDB

test:
	go test -v -race -cover ./...