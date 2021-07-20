#@IgnoreInspection BashAddShebang

export APP=gokkan

export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))

export POSTGRES_ADDRESS="localhost:5432"
export LDFLAGS="-w -s"
export POSTGRES_DSN="postgres://gokkan:1@$(POSTGRES_ADDRESS)/gokkan?sslmode=disable"
all: format lint build

############################################################
# Build and Run
############################################################

build:
	CGO_ENABLED=1 go build -ldflags $(LDFLAGS)  ./cmd/gokkan

install:
	CGO_ENABLED=1 go install -ldflags $(LDFLAGS) ./cmd/gokkan

############################################################
# Format and Lint
############################################################

check-formatter:
	which goimports || GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports
	which gci || GO111MODULE=off go get github.com/daixiang0/gci
	which gofumpt || GO111MODULE=off go get mvdan.cc/gofumpt

format: check-formatter
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R goimports -w R
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R gofmt -s -w R
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R gci -w R
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R gofumpt -s -w R


check-linter:
	which golangci-lint || GO111MODULE=off curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sudo sh -s -- -b $(go env GOPATH)/bin v1.40.1

lint: check-linter format
	golangci-lint run $(ROOT)/...

############################################################
# Test
############################################################

test:
	go test -v -race -p 1 ./...

############################################################
# Install git hook
############################################################

install-hook:
	git config --local core.hooksPath ./githooks


############################################################
# docker-compose
############################################################

up:
	docker-compose up -d
down:
	docker-compose down
ps: up
	docker-compose ps

############################################################
# SQL-Migration
############################################################

check-migrate:
	which migrate || go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-create:
	migrate create -ext sql -dir ./config/sql $(NAME)

migrate-up:
	migrate -verbose  -path ./config/sql -database $(POSTGRES_DSN) up

migrate-down:
	 migrate -path ./config/sql -database $(POSTGRES_DSN) down

migrate-reset:
	 migrate -path ./config/sql -database $(POSTGRES_DSN) drop