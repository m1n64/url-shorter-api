# Variables
SERVICES = infrastructure-service user-service
SHARED_NETWORK_DELAY = 15

.PHONY: help
help:
	@echo "Usage:"
	@echo "  make up - Start all services"
	@echo "  make stop - Stop all services"
	@echo "  make restart - Restart all services"
	@echo "  make network - Create docker network"
	@echo "  make up-service SERVICE= - Start service"
	@echo "  make stop-service SERVICE= - Stop service"
	@echo "  make restart-service SERVICE= - Restart service"

.PHONY: init
init:
	@echo "🛠️  Initialization all services and networks..."
	$(MAKE) network
	$(MAKE) up

.PHONY: up
up:
	@echo "🚀 Start all serviceS..."
	@for service in $(SERVICES); do \
    	if [ "$$service" = "infrastructure-service" ]; then \
    		echo "🟢 Starting $$service with default compose file..."; \
    		(cd $$service && docker-compose up -d || echo "❌ Error starting $$service"); \
    		echo "⏳ Waiting for shared-network to be ready for $(SHARED_NETWORK_DELAY) seconds..."; \
            sleep $(SHARED_NETWORK_DELAY); \
    	else \
    		echo "🟢 Starting $$service with dev compose file..."; \
    		(cd $$service && docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d || echo "❌ Error starting $$service"); \
    	fi; \
    done

.PHONY: stop
stop:
	@echo "📴 Stop all services..."
	@for service in $(SERVICES); do \
		echo "🔴 Stopping $$service..."; \
		(cd $$service && docker-compose stop || echo "❌ Error stop $$service"); \
	done

.PHONY: restart
restart:
	@echo "🔄 Reload all services..."
	@for service in $(SERVICES); do \
  		if [ "$$service" = "infrastructure-service" ]; then \
      		echo "🔄 Reload $$service with default compose file..."; \
      		(cd $$service && docker-compose restart || echo "❌ Error restarting $$service"); \
      		echo "⏳ Waiting for shared-network to be ready for $(SHARED_NETWORK_DELAY) seconds..."; \
            sleep $(SHARED_NETWORK_DELAY); \
      	else \
      		echo "🔄 Reload $$service with dev compose file..."; \
      		(cd $$service && docker-compose -f docker-compose.yml -f docker-compose.dev.yml restart || echo "❌ Error restarting $$service"); \
      	fi; \
	done

.PHONY: restart-service
restart-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "❌ Please specify a service to restart, e.g., 'make restart-service SERVICE=user-service'"; \
		exit 1; \
	fi
	@echo "🔄 Restarting service: $(SERVICE)..."
	@if [ "$(SERVICE)" = "infrastructure-service" ]; then \
    	(cd $(SERVICE) && docker-compose up -d || echo "❌ Error starting $(SERVICE)"); \
    else \
    	(cd $(SERVICE) && docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d || echo "❌ Error restarting $(SERVICE)"); \
    fi

.PHONY: up-service
up-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "❌ Please specify a service to start, e.g., 'make up-service SERVICE=user-service'"; \
		exit 1; \
	fi
	@echo "🚀 Starting service: $(SERVICE)..."
	@if [ "$(SERVICE)" = "infrastructure-service" ]; then \
    	(cd $(SERVICE) && docker-compose up -d || echo "❌ Error starting $(SERVICE)"); \
    else \
    	(cd $(SERVICE) && docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d || echo "❌ Error starting $(SERVICE)"); \
    fi

.PHONY: stop-service
stop-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "❌ Please specify a service to stop, e.g., 'make stop-service SERVICE=user-service'"; \
		exit 1; \
	fi
	@echo "📴 Stopping service: $(SERVICE)..."
	@(cd $(SERVICE) && docker-compose stop || echo "❌ Error stopping $(SERVICE)")

.PHONY: network
network:
	@echo "🌐 Create network..."
	docker network create tidy-url-network || echo "❌ Network \"tidy-url-network\" already exists"