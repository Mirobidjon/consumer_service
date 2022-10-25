FROM golang:1.17 as builder

#
RUN mkdir -p $GOPATH/src/gitlab-dev.soliqservis.uz/task/consumer_service
WORKDIR $GOPATH/src/gitlab-dev.soliqservis.uz/task/consumer_service

# Copy the local package files to the container's workspace.
COPY . ./

# installing depends and build
RUN export CGO_ENABLED=0 && \
    export GOOS=linux && \
    go mod vendor && \
    make build && \
    mv ./bin/consumer_service /

FROM alpine
COPY --from=builder consumer_service .

ENTRYPOINT ["/consumer_service"]
