PRICE_STATS_BINARY=price-stats-service
PRODUCT_HISTORY_BINARY=product-history-service 


## up: stops docker-compose (if running), builds all projects and starts docker compose
up:
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"
	
## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

build-all:
	$(MAKE) -j2 build-price-stats build-product-history build-docker

test-all:
	$(MAKE) -j2 test-price-stats test-product-history

build-docker:
	docker-compose build --no-cache product-history price-stats

## build-price-stats: builds the price-stats-service binary as a linux executable
build-price-stats:
	@echo "Building price-stats-service binary..."
	cd ./price-stats-service && env GOOS=linux CGO_ENABLED=0 go build -o ${PRICE_STATS_BINARY} ./cmd
	@echo "Done!"


## build-product-history: builds the product-history-service binary as a linux executable
build-product-history:
	@echo "Building product-history-service binary..."
	cd ./product-history-service  && env GOOS=linux go build -o ${PRODUCT_HISTORY_BINARY} ./cmd
	@echo "Done!"

test-price-stats:
	cd price-stats-service && go clean -testcache && go test -race -v 2>&1 ./...
test-product-history:
	cd product-history-service && go clean -testcache && go test -race -v 2>&1 ./...


