BACKEND_BINARY=backend_binary

## UP: Starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## DOWN: Stop docker containers
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## BUILDS: Builds the binaries as a linux executables for our docker images
build:
	@echo "Building backend binary..."
	cd ./backend && env GOOS=linux CGO_ENABLED=0 go build -o ${BACKEND_BINARY} .
	@echo "Done!"
