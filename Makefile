.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help message
	@echo "Available commands:"
	@echo ""
	@echo "Backend:"
	@echo "  make backend-dev              - Run a dev backend server"
	@echo "  make backend-build            - Build the backend server"
	@echo "  make backend-format           - Format backend code"
	@echo "  make backend-test             - Run backend tests"
	@echo "  make backend-test-verbose     - Run backend tests with verbose output"
	@echo "  make backend-test-coverage    - Run backend tests with coverage"
	@echo ""
	@echo "Frontend:"
	@echo "  make frontend-dev             - Run a dev frontend server"
	@echo "  make frontend-lint            - Lint the frontend code"
	@echo "  make frontend-lint-fix        - Run the linter and fix fixable issues"
	@echo "  make frontend-format          - Run prettifier on the frontend code"
	@echo "  make frontend-typecheck       - Run typecheck on the frontend code"
	@echo "  make frontend-test            - Run frontend tests"
	@echo "  make frontend-test-watch      - Run frontend tests in watch mode"
	@echo "  make frontend-build           - Build the frontend code"
	@echo ""
	@echo "Infrastructure:"
	@echo "  make infra-start              - Start the development docker infrastructure"
	@echo "  make infra-stop               - Stop the development docker infrastructure"
	@echo ""
	@echo "Stripe:"
	@echo "  make stripe-listen            - Start the stripe webhook event listener"

# Catch-all target to prevent Make from complaining about unknown targets
%:
	@:

backend-dev:
	cd backend; gow run server.go

backend-build:
	cd backend; go build -o bin/server server.go

backend-format:
	cd backend; go fmt ./...

backend-test:
	cd backend; TEST_MODE=true go test ./...

backend-test-verbose:
	cd backend; TEST_MODE=true go test -v ./...

backend-test-coverage:
	cd backend; TEST_MODE=true go test -cover ./...

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

frontend-test:
	cd frontend; npm test

frontend-test-watch:
	cd frontend; npm test:watch

frontend-build:
	cd frontend; npm run build

infra-start:
	docker compose up -d

infra-stop:
	docker compose down

stripe-listen:
	stripe listen \
		-H "Content-Type: application/json" \
		--forward-to localhost:4040/api/webhooks/stripe/account/snapshot \
		--forward-connect-to localhost:4040/api/webhooks/stripe/account/snapshot \
		--events "invoice.paid,invoice.payment_failed,invoice.payment_action_required,customer.subscription.created,customer.subscription.updated,customer.subscription.deleted" \
		--forward-thin-to localhost:4040/api/webhooks/stripe/connect/thin \
		--forward-thin-connect-to localhost:4040/api/webhooks/stripe/connect/thin \
		--thin-events "v2.core.account.created,v2.core.account.updated,v2.core.account.closed,v2.core.account_person.updated,v2.core.account[identity].updated,v2.core.account[configuration.customer].capability_status_updated,v2.core.account[configuration.merchant].capability_status_updated,v2.core.account[configuration.recipient].capability_status_updated,v2.core.account[requirements].updated,v2.core.account.updated" \
