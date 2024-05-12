GOLANGCI_LINT_CACHE?=/tmp/praktikum-golangci-lint-cache
CONFIG_FILE=not_set
MASTER_KEY=change_me

.PHONY: golangci-lint-run
golangci-lint-run: _golangci-lint-rm-unformatted-report

.PHONY: _golangci-lint-reports-mkdir
_golangci-lint-reports-mkdir:
	mkdir -p ./golangci-lint

.PHONY: _golangci-lint-run
_golangci-lint-run: _golangci-lint-reports-mkdir
	-docker run --rm \
    -v $(shell pwd):/app \
    -v $(GOLANGCI_LINT_CACHE):/root/.cache \
    -w /app \
    golangci/golangci-lint:v1.55.2 \
        golangci-lint run \
            -c .golangci.yml \
	> ./golangci-lint/report-unformatted.json

.PHONY: _golangci-lint-format-report
_golangci-lint-format-report: _golangci-lint-run
	cat ./golangci-lint/report-unformatted.json | jq > ./golangci-lint/report.json

.PHONY: _golangci-lint-rm-unformatted-report
_golangci-lint-rm-unformatted-report: _golangci-lint-format-report
	rm ./golangci-lint/report-unformatted.json

.PHONY: golangci-lint-clean
golangci-lint-clean:
	sudo rm -rf ./golangci-lint

.PHONY: build
build: _create_build_dir
	go build -o ./build/gopher ./cmd/gophermart
	chmod +x ./build/gopher

.PHONY: _create_build_dir
_create_build_dir:
	mkdir -p ./build

.PHONY: ggen
ggen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/pass.proto

.PHONY: start-server
start-server:
	go build -ldflags "-X main.configFilePath=$(CONFIG_FILE) -X main.masterKey=$(MASTER_KEY)" -o build/temp_server cmd/server/main.go
	-./build/temp_server

.PHONY: bc
bc:
	go build -o build/passkeep cmd/client/main.go