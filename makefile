.PHONY: all build up down run test swagger

# Application and database container names
APP_NAME=app

# Default target: build and start the application
all: build up

# Build the application
build:
	@echo "Building the Go application..."
	docker-compose build $(APP_NAME)

# Start the containers
up:
	@echo "Starting the application..."
	docker-compose up -d $(APP_NAME)

# Stop the containers
down:
	@echo "Stopping all containers..."
	docker-compose down

# Run the application
run:
	docker-compose -f docker-compose.yaml up -d
##
swagger:
	export GOFLAGS="-mod=mod" && swag init --output ./docs --generalInfo ./cmd/main.go --parseInternal --parseDependency && export GOFLAGS="-mod=vendor"
test:
	@echo "Running tests..."
	go test ./integration_test/... -v -cover
