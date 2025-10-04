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

stripe-listen:
	stripe listen \
		-H "Content-Type: application/json" \
		--forward-to localhost:4040/api/webhooks/stripe \
		--forward-connect-to localhost:4040/api/webhooks/stripe \
		--forward-thin-to localhost:4040/api/webhooks/stripe \
		--forward-thin-connect-to localhost:4040/api/webhooks/stripe \
		--thin-events "v2.core.account.updated,v2.core.account.closed,v2.core.account[configuration.customer].capability_status_updated,v2.core.account[configuration.merchant].capability_status_updated,v2.core.account[configuration.recipient].capability_status_updated,v2.core.account[requirements].updated" \
		--events "capability.updated,account.updated"
