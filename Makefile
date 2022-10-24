CURRENT_DIR=$(shell pwd)

APP=$(shell basename ${CURRENT_DIR})
APP_CMD_DIR=${CURRENT_DIR}/cmd

TAG=latest
ENV_TAG=latest

-include .env

POSTGRESQL_URL='postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable'

pull-proto-module:
	git submodule update --init --recursive

update-proto-module:
	git submodule update --remote --merge

copy-proto-module:
	rm -rf ${CURRENT_DIR}/protos
	rsync -rv --exclude={'/.git','LICENSE','README.md'} ${CURRENT_DIR}/ss_protos/* ${CURRENT_DIR}/protos

gen-proto-module:
	./scripts/gen_proto.sh ${CURRENT_DIR}
migrate-local-up:
	migrate -database ${POSTGRESQL_URL} -path migrations/postgres up

migrate-local-down:
	migrate -database ${POSTGRESQL_URL} -path migrations down

migration-up:
	migrate -path ./migrations/postgres -database 'postgres://postgres:123@0.0.0.0:5432/consumer_service?sslmode=disable' up

migration-down:
	migrate -path ./migrations/postgres -database 'postgres://postgres:123@0.0.0.0:5432/consumer_service?sslmode=disable' down

build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

build-image:
	docker build --rm -t ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG} .
	docker tag ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG} ${REGISTRY}/${PROJECT_NAME}/${APP}:${ENV_TAG}

push-image:
	docker push ${REGISTRY}/${PROJECT_NAME}/${APP}:${TAG}
	docker push ${REGISTRY}/${PROJECT_NAME}/${APP}:${ENV_TAG}

swag-init:
	swag init -g api/api.go -o api/docs

run:
	go run cmd/main.go

linter:
	golangci-lint run

push:
	make build-image TAG=1 SERVICE_NAME=consumer_service PROJECT_NAME=learn-cloud-0809 REGISTRY=us.gcr.io ENV_TAG=latest 
	make push-image TAG=1 SERVICE_NAME=consumer_service PROJECT_NAME=learn-cloud-0809 REGISTRY=us.gcr.io ENV_TAG=latest 
