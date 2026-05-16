.PHONY: build run dev format lint lint-fix

# Build the Wails application for production
build:
	wails build

# Build and run the production binary
run: build
	./build/bin/switchy.exe

# Run the application in development mode with hot-reload
dev:
	wails dev

# Format Go code and run Prettier on the frontend
format:
	go fmt ./...
	cd frontend && npm run format

# Lint Go code and frontend (check only)
lint:
	go vet ./...
	cd frontend && npm run lint

# Lint and auto-fix frontend issues
lint-fix:
	cd frontend && npm run lint:fix
