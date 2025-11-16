.DEFAULT_GOAL := help

# List of known targets to filter out from arguments
KNOWN_TARGETS := help backend-dev backend-build backend-format backend-test backend-test-verbose backend-test-coverage \
                 frontend-dev frontend-lint frontend-lint-fix frontend-format frontend-typecheck frontend-test frontend-test-watch frontend-build \
                 infra-start infra-stop stripe-listen

# Variable to capture additional arguments
# Usage: make backend-test ARGS=./internal/organizations
# Or: make backend-test ./internal/organizations
# If ARGS is not explicitly set, capture positional arguments
ifndef ARGS
  ARGS := $(filter-out $(KNOWN_TARGETS),$(MAKECMDGOALS))
endif

# Normalize paths: ensure relative paths start with ./
# This handles cases where Make might strip the ./ prefix
normalize-path = $(if $(filter /% http://% https://%,$(1)),$(1),$(if $(filter ./%,$(1)),$(1),./$(1)))

# Normalize ARGS if it's set
NORMALIZED_ARGS := $(if $(ARGS),$(call normalize-path,$(ARGS)),)

# Default ARGS for test commands if not provided
TEST_ARGS := $(if $(NORMALIZED_ARGS),$(NORMALIZED_ARGS),./...)

.PHONY: help
help: ## Show this help message
	@echo "Available commands:"
	@echo ""
	@echo "Backend:"
	@echo "  make backend-dev              - Run a dev backend server"
	@echo "  make backend-build            - Build the backend server"
	@echo "  make backend-format           - Format backend code"
	@echo "  make backend-test [ARGS=...]  - Run backend tests (e.g., make backend-test ./internal/organizations)"
	@echo "  make backend-test-verbose [ARGS=...]  - Run backend tests with verbose output"
	@echo "  make backend-test-coverage [ARGS=...]  - Run backend tests with coverage"
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
	cd backend; go fmt $(if $(NORMALIZED_ARGS),$(NORMALIZED_ARGS),./...)

backend-test:
	cd backend; TEST_MODE=true go test $(TEST_ARGS)

backend-test-verbose:
	cd backend; TEST_MODE=true go test -v $(TEST_ARGS)

backend-test-coverage:
	cd backend; TEST_MODE=true go test -cover $(TEST_ARGS)

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
