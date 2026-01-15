BACKEND_BINARY=backend_binary

## START: Starts the frontend and backend
start:
	@echo "Starting Commandup..."
	podman compose up -d postgres &
	cd ./frontend && pnpm start &
	cd ./backend && air &
	@echo "Done!"

## START_DB: Start database container only
start_database:
	@echo "Starting database..."
	podman compose up -d postgres
	@echo "Database started!"

## STOP_DB: Stop database container only
stop_database:
	@echo "Stopping database..."
	podman compose stop postgres
	@echo "Database stopped!"


## DOWN: Stop docker containers
down:
	@echo "Stopping docker compose..."
	podman compose down
	@echo "Done!"

## BUILDS: Builds the binaries as a linux executables for our docker images
build:
	@echo "Building backend binary..."
	cd ./backend && env GOOS=linux CGO_ENABLED=0 go build -o ${BACKEND_BINARY} .
	@echo "Done!"
