.PHONY: backend frontend dev

backend:
	go run ./cmd/api/main.go


frontend:
	cd frontend && npm run dev

dev:
	go run ./cmd/api/main.go &
	cd frontend && npm run dev
