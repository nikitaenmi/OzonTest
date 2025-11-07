PROTO_DIR=api
GEN_DIR=gen
THIRD_PARTY_DIR=third_party
BIN_DIR=bin

.PHONY: start
start: check-env deps gen docker-up

.PHONY: stop
stop: docker-down

.PHONY: restart
restart: stop start

.PHONY: clean
clean: docker-down
	rm -rf $(GEN_DIR)/* $(BIN_DIR)/* .env

.PHONY: test-load
test-load:
	cd loadtest && go run main.go

.PHONY: check-env
check-env:
	@if [ ! -f .env ]; then \
		cp samples/.env.example .env; \
	fi

.PHONY: deps
deps:
	go mod download
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.15.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.15.0

.PHONY: gen
gen: check-env
	protoc \
		--proto_path=$(PROTO_DIR) \
		--proto_path=$(THIRD_PARTY_DIR) \
		--proto_path=$(THIRD_PARTY_DIR)/protoc-gen-openapiv2 \
		--go_out=$(GEN_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(GEN_DIR) \
		--grpc-gateway_opt=paths=source_relative \
		--grpc-gateway_opt=logtostderr=true \
		--openapiv2_out=$(GEN_DIR) \
		--openapiv2_opt=logtostderr=true \
		$(PROTO_DIR)/payment.proto

.PHONY: docker-up
docker-up: check-env
	docker-compose up -d
	@sleep 5
	@until docker-compose exec dbstore pg_isready -U root; do \
		sleep 2; \
	done

.PHONY: docker-down
docker-down:
	docker-compose down

.PHONY: build
build: gen
	go build -o $(BIN_DIR)/server cmd/main.go

.PHONY: run
run: build
	./$(BIN_DIR)/server