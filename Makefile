backend-dev:
	cd backend; gow run server.go

backend-build:
	cd backend; go build -o bin/server server.go

frontend-dev:
	cd frontend; npm run dev 

frontend-lint:
	cd frontend; npm run lint

frontend-lint-fix:
	cd frontend; npm run lint:fix

frontend-format:
	cd frontend; npm run format

frontend-typecheck:
	cd frontend; npm run check

frontend-build:
	cd frontend; npm run build

infra-start:
	docker compose up -d

infra-stop:
	docker compose down

stripe-listen:
	stripe listen \
		-H "Content-Type: application/json" \
		--forward-to localhost:4040/api/webhooks/stripe/snapshot \
		--forward-connect-to localhost:4040/api/webhooks/stripe/snapshot \
		--events "invoice.paid,invoice.payment_failed,invoice.payment_action_required" \
		--forward-thin-to localhost:4040/api/webhooks/stripe/thin \
		--forward-thin-connect-to localhost:4040/api/webhooks/stripe/thin \
		--thin-events "v2.core.account.updated,v2.core.account.closed,v2.core.account[configuration.customer].capability_status_updated,v2.core.account[configuration.merchant].capability_status_updated,v2.core.account[configuration.recipient].capability_status_updated,v2.core.account[requirements].updated" \
