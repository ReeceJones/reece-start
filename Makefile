backend-dev:
	cd backend; gow run server.go

frontend-dev:
	cd frontend; npm run dev 

infra-start:
	docker compose up -d

infra-stop:
	docker compose down
