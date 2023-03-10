launch_args=
test_args=-coverprofile cover.out && go tool cover -func cover.out
cover_args=-cover -coverprofile=cover.out `go list ./...` && go tool cover -html=cover.out

SERVICE_NAME=storage-service
VERSION?=dev
DOCKER_IMAGE_NAME=krobus00/${SERVICE_NAME}
CONFIG?=./config.yml
NAMESPACE?=default

# make tidy
tidy:
	go mod tidy

# make proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. \
  		--go-grpc_opt=paths=source_relative pb/storage/*.proto
	ls pb/storage/*.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'

# make lint
lint:
	golangci-lint run --disable-all -E errcheck -E misspell -E revive -E goimports

# make run-dev server, make run-dev worker
run-dev:
ifeq (server, $(filter server,$(MAKECMDGOALS)))
	$(eval launch_args=server $(launch_args))
else ifeq (worker, $(filter worker,$(MAKECMDGOALS)))
	$(eval launch_args=worker $(launch_args))
endif
	air --build.cmd "go build -o bin/storage-service main.go" --build.bin "./bin/storage-service $(launch_args)"
	
# make build
build:
	# build binary file
	go build -ldflags "-s -w" -o ./bin/storage-service ./main.go
ifeq (, $(shell which upx))
	$(warning "upx not installed")
else
	# compress binary file if upx command exist
	upx -9 ./bin/storage-service
endif

# make image VERSION="vx.x.x"
image:
	docker build -t ${DOCKER_IMAGE_NAME}:${VERSION} . -f ./deployments/Dockerfile --build-arg GITHUB_USERNAME=${GITHUB_USERNAME} --build-arg GITHUB_TOKEN=${GITHUB_TOKEN}

# make deploy VERSION="vx.x.x"
# make deploy VERSION="vx.x.x" NAMESPACE="staging"
# make deploy VERSION="vx.x.x" NAMESPACE="staging" CONFIG="./config-staging.yml"
deploy:
	helm upgrade --install storage-service ./deployments/helm/server-storage-service --set-file appConfig="${CONFIG}" --set app.container.version="${VERSION}" -n ${NAMESPACE}

# make test
test:
ifeq (, $(shell which richgo))
	go test ./... $(test_args)
else
	richgo test ./... $(test_args)
endif

# make cover
cover: test
ifeq (, $(shell which richgo))
	go test $(cover_args)
else
	richgo test $(cover_args)
endif

%:
	@: