#!/bin/bash

# Start the Go backend
cd backend
go run main.go &
cd ..

# Start the React frontend
cd frontend
pnpm start
