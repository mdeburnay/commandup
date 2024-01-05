#!/bin/bash

# Start the Go backend
cd backend
$(go env GOPATH)/bin/air &
cd ..

# Start the React frontend
cd frontend
pnpm start
