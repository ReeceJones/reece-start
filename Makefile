backend-dev:
	cd backend; gow run server.go

backend:
	cd backend; go build -o bin/server server.go

frontend-dev:
	cd frontend; npm run dev 

infra-start:
	docker compose up -d

infra-stop:
	docker compose down
