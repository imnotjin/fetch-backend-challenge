.PHONY: setup up down test

setup:
	@echo "Creating .env and .env.test files..."
	@echo "DB_USER=postgres" > .env
	@echo "DB_PASSWORD=your_secure_password" >> .env
	@echo "DB_NAME=points_db" >> .env
	@echo "DB_HOST=db" >> .env
	@echo "DB_PORT=5432" >> .env
	@echo "POSTGRES_USER=postgres" >> .env
	@echo "POSTGRES_PASSWORD=your_secure_password" >> .env
	@echo "POSTGRES_DB=points_db" >> .env

	@echo "DB_USER=postgres" > .env.test
	@echo "DB_PASSWORD=your_secure_password" >> .env.test
	@echo "DB_NAME=test_db" >> .env.test
	@echo "DB_HOST=localhost" >> .env.test
	@echo "DB_PORT=5433" >> .env.test
	@echo "POSTGRES_USER=postgres" >> .env.test
	@echo "POSTGRES_PASSWORD=your_secure_password" >> .env.test
	@echo "POSTGRES_DB=test_db" >> .env.test
	@echo ".env and .env.test files created."

up:
	@docker-compose up --build

down:
	@docker-compose down

test:
	@go test ./handlers -v
