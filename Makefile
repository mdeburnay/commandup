BACKEND_BINARY=backend_binary

## START: Starts the frontend, backend, and database using a Procfile runner (e.g. hivemind)
start:
	@if command -v hivemind >/dev/null 2>&1; then \
		hivemind; \
	elif command -v overmind >/dev/null 2>&1; then \
		overmind start; \
	elif command -v foreman >/dev/null 2>&1; then \
		foreman start; \
	else \
		echo "No Procfile runner (hivemind, overmind, or foreman) found."; \
		echo "Please install hivemind: brew install hivemind"; \
		echo "Or run: make start_legacy"; \
	fi

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
