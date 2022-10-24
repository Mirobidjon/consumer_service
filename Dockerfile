FROM golang:1.17 as builder

#
RUN mkdir -p $GOPATH/src/gitlab-dev.soliqservis.uz/task/consumer_service
WORKDIR $GOPATH/src/gitlab-dev.soliqservis.uz/task/consumer_service
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.43.0

# Copy the local package files to the container's workspace.
COPY . ./

# installing depends and build
RUN export CGO_ENABLED=0 && \
    export GOOS=linux && \
    go mod vendor && \
    make linter && \
    make build && \
    mv ./bin/consumer_service /

FROM alpine
COPY --from=builder consumer_service .
RUN apk update && apk add -U tzdata && cp /usr/share/zoneinfo/Asia/Tashkent /etc/localtime && apk del tzdata

ENTRYPOINT ["/consumer_service"]
