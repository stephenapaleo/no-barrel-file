# ==================================================================================== #
# VARIABLES
# ==================================================================================== #

APP_NAME = no-barrel-file
DOCKER_REGISTRY = nergie42
DOCKER_IMAGE = ${APP_NAME}
DOCKER_IMAGE_WITH_TAG = ${DOCKER_IMAGE}:${DOCKER_TAG}
BUILD_DIR = ./bin


# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

.PHONY: confirm
confirm:
  # Use printf instead of echo because the builtin echo in sh(1) does not accept the -n option. see man echo and https://unix.stackexchange.com/questions/700675/why-is-echo-e-behaving-weird-in-a-makefile for more information.
	@printf "%s" 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## run/display: display barrel files
.PHONY: run/display
run/display:
	@echo "Displaying barrel files..."
	go run . display --root-path ./tests/data/input

## run/count: count barrel files
.PHONY: run/count
run/count:
	@echo "Counting barrel files using..."
	go run . count --root-path ./tests/data/input

## run/replace: replace barrel files imports
.PHONY: run/replace
run/replace:
	@echo "Updating barrel files imports..."
	cp -rf tests/data/input /tmp/
	@echo "Updating barrel files..."
	go run . replace --root-path /tmp/input --ignore-paths ignored --alias-config-path tsconfig.json -v


# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## audit: tidy and vendor dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor


# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build: build the cli application
.PHONY: build
build:
	@echo "Building ${APP_NAME}..."
	go build -ldflags='-s' -o=${BUILD_DIR}/${APP_NAME} .
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=${BUILD_DIR}/linux-amd64/${APP_NAME} .

## clean: clean up build artifacts
.PHONY: clean
clean: confirm
	@echo "Build directory cleaning up..."
	rm -rf ${BUILD_DIR}


# ==================================================================================== #
# DOCKER
# ==================================================================================== #

## docker-build: build the Docker image
.PHONY: docker-build
docker-build:
	@echo "Building Docker image ${DOCKER_IMAGE_WITH_TAG}..."
	docker build -t ${DOCKER_REGISTRY}/${DOCKER_IMAGE_WITH_TAG} .

## docker-clean: remove the Docker image
.PHONY: docker-clean
docker-clean: confirm
	@echo "Removing Docker image ${DOCKER_IMAGE_WITH_TAG}..."
	docker rmi -f ${DOCKER_REGISTRY}/${DOCKER_IMAGE_WITH_TAG} || true

## docker-tag-latest: add tag latest to the Docker image
.PHONY: docker-tag-latest
docker-tag-latest:
	@echo "Tagging Docker image ${DOCKER_REGISTRY}/${DOCKER_IMAGE_WITH_TAG} with tag latest..."
	docker tag ${DOCKER_REGISTRY}/${DOCKER_IMAGE_WITH_TAG} ${DOCKER_REGISTRY}/${DOCKER_IMAGE}:latest

## docker-push: push the Docker image to a registry
.PHONY: docker-push
docker-push:
	@echo "Pushing Docker image ${DOCKER_IMAGE_WITH_TAG} to registry ${DOCKER_REGISTRY}..."
	docker push ${DOCKER_REGISTRY}/${DOCKER_IMAGE_WITH_TAG}

## docker-run: run the Docker container
.PHONY: docker-run
docker-run:
	@echo "Running ${APP_NAME} in Docker with image ${DOCKER_REGISTRY}/${DOCKER_IMAGE_WITH_TAG}..."
	docker run --rm ${DOCKER_REGISTRY}/${DOCKER_IMAGE_WITH_TAG} --help
