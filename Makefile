launch_args=
test_args=-coverprofile cover.out && go tool cover -func cover.out
cover_args=-cover -coverprofile=cover.out `go list ./...` && go tool cover -html=cover.out

SERVICE_NAME=storage-service
VERSION?= $(shell git describe --match 'v[0-9]*' --tags --always)
DOCKER_IMAGE_NAME=krobus00/${SERVICE_NAME}
CONFIG?=./config.yml
NAMESPACE?=default
PROJECT_REPO=github.com/krobus00/${SERVICE_NAME}
MIGRATION_ACTION?="up"
MIGRATION_NAME?=""
MIGRATION_STEP?="999"

# make tidy
tidy:
	go mod tidy

# make clean-up-mock
clean-up-mock:
	rm -rf ./internal/model/mock
	rm -rf ./pb/storage/mock


pb/storage/mock/mock_storage_service_client.go:
	mockgen -destination=pb/storage/mock/mock_storage_service_client.go -package=mock ${PROJECT_REPO}/pb/storage StorageServiceClient

# make generate
generate: clean-up-mock pb/storage/mock/mock_storage_service_client.go
	go generate ./...


# make proto
proto:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. \
  		--go-grpc_opt=paths=source_relative pb/storage/*.proto
	@ls pb/storage/*.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'

# make lint
lint:
	@golangci-lint run

# make run dev server
# make run dev worker
# make run server
# make run worker
# make run migration
# make run migration MIGRATION_ACTION=up
# make run migration MIGRATION_ACTION=create MIGRATION_NAME=create_table_products
# make run migration MIGRATION_ACTION=up MIGRATION_STEP=1
run:
ifeq (dev server, $(filter dev server,$(MAKECMDGOALS)))
	$(eval launch_args=server $(launch_args))
	air --build.cmd 'go build -ldflags "-s -w -X main.version=$(VERSION) -X main.name=$(SERVICE_NAME)" -o bin/storage-service main.go' --build.bin "./bin/storage-service $(launch_args)"
else ifeq (dev worker, $(filter dev worker,$(MAKECMDGOALS)))
	$(eval launch_args=worker $(launch_args))
	air --build.cmd 'go build -ldflags "-s -w -X main.version=$(VERSION) -X main.name=$(SERVICE_NAME)" -o bin/storage-service main.go' --build.bin "./bin/storage-service $(launch_args)"
else ifeq (worker, $(filter worker,$(MAKECMDGOALS)))
	$(eval launch_args=worker $(launch_args))
	$(shell if test -s ./bin/storage-service; then ./bin/storage-service $(launch_args); else echo storage binary not found; fi)
else ifeq (server, $(filter server,$(MAKECMDGOALS)))
	$(eval launch_args=server $(launch_args))
	$(shell if test -s ./bin/storage-service; then ./bin/storage-service $(launch_args); else echo storage binary not found; fi)
else ifeq (migration, $(filter migration,$(MAKECMDGOALS)))
	$(shell if ! test -s ./bin/auth-service; then go build -ldflags "-s -w -X main.version=$(VERSION) -X main.name=$(SERVICE_NAME)"  -o ./bin/storage-service ./main.go; fi)
	$(eval launch_args=migration --action $(MIGRATION_ACTION) --name $(MIGRATION_NAME) --step $(MIGRATION_STEP) $(launch_args))
	./bin/storage-service $(launch_args)
endif

# make build
build:
	# build binary file
	go build -ldflags "-s -w -X main.version=$(VERSION) -X main.name=$(SERVICE_NAME)" -o ./bin/storage-service ./main.go
ifeq (, $(shell which upx))
	$(warning "upx not installed")
else
	# compress binary file if upx command exist
	upx -9 ./bin/storage-service
endif

# make image VERSION="vx.x.x"
image:
	docker build -t ${DOCKER_IMAGE_NAME}:${VERSION} . -f ./deployments/Dockerfile

# make push-image VERSION="vx.x.x"
push-image:
	docker push ${DOCKER_IMAGE_NAME}:${VERSION}

# make docker-build-push VERSION="vx.x.x"
docker-build-push: image push-image

# make deploy VERSION="vx.x.x"
# make deploy VERSION="vx.x.x" NAMESPACE="staging"
# make deploy VERSION="vx.x.x" NAMESPACE="staging" CONFIG="./config-staging.yml"
deploy:
	@helm upgrade --install storage-service ./deployments/helm/server-storage-service --set-file appConfig="${CONFIG}" --set app.container.version="${VERSION}" -n ${NAMESPACE}

# make coverage
coverage:
	@echo "total code coverage : "
	@go tool cover -func cover.out | grep total | awk '{print substr($$3, 1, length($$3)-1)}'

# make test
test:
ifeq (, $(shell which richgo))
	go test ./... $(test_args)
else
	richgo test ./... $(test_args)
endif

# make cover
cover:
ifeq (, $(shell which richgo))
	go test $(cover_args)
else
	richgo test $(cover_args)
endif

# make changelog VERSION=vx.x.x
changelog:
	git-chglog -o CHANGELOG.md --next-tag $(VERSION)

%:
	@:
