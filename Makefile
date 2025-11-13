.PHONY: help setup-backend setup-frontend run-backend run-frontend run-docker clean

help:
	@echo "Available commands:"
	@echo "  make setup-backend    - Setup Go backend dependencies"
	@echo "  make setup-frontend   - Setup frontend dependencies"
	@echo "  make run-backend      - Run backend server"
	@echo "  make run-frontend     - Run frontend dev server"
	@echo "  make run-docker       - Run with Docker Compose"
	@echo "  make clean            - Clean generated files"

setup-backend:
	cd backend && go mod tidy && go mod download

setup-frontend:
	cd frontend && npm install

run-backend:
	cd backend && go run main.go

run-frontend:
	cd frontend && npm run dev

run-docker:
	docker-compose up --build

clean:
	rm -rf backend/data backend/uploads backend/generated
	rm -rf frontend/.next frontend/node_modules
	rm -rf backend/vendor


