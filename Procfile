# Procfile for Commandup
# Use a tool like 'hivemind' or 'overmind' to run this.

# Frontend (React)
web: cd frontend && PORT=3000 pnpm start

# Backend (Go with Air for hot-reloading)
api: cd backend && air

# Database (Postgres via Podman Compose)
# We run this in the foreground so the Procfile manager can track it.
db: podman compose up postgres
