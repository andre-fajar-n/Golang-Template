.PHONY: generate build run doc validate spec clean help

all: clean spec generate build

validate:
	swagger validate ./api/go-template/index.yml

spec:
	swagger expand --output=./api/go-template/result.yml --format=yaml ./api/go-template/index.yml

build: 
	CGO_ENABLED=0 GOOS=linux go build -v -installsuffix cgo ./cmd/server
	
run:
	./server --port=8080 --host=0.0.0.0 --config=./configs/app.yaml

run-local:
	go run cmd/server/main.go --port=8080

doc: validate
	swagger serve api/go-template/index.yml --no-open --host=0.0.0.0 --port=8080 --base-path=/

clean:
	rm -rf server
	rm -rf ./gen/models
	rm -rf ./gen/rest
	go clean -i .

generate: validate
	swagger generate server --exclude-main -A server -t gen -f ./api/go-template/result.yml --principal models.Principal

help:
	@echo "make: compile packages and dependencies"
	@echo "make validate: OpenAPI validation"
	@echo "make spec: OpenAPI Spec"
	@echo "make clean: remove object files and cached files"
	@echo "make build: Generate Server and Client API"
	@echo "make doc: Serve the Doc UI"
	@echo "make run: Serve binary file"
	@echo "make run-local: Serve main.go"