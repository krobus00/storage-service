FROM golang:1.18.3-alpine as base

ARG GITHUB_USERNAME
ARG GITHUB_TOKEN

WORKDIR /builder
# install upx and make
RUN apk add upx make
ENV GO111MODULE=on CGO_ENABLED=0
RUN apk add git
RUN git config --global url."https://${GITHUB_USERNAME}:${GITHUB_TOKEN}@github.com".insteadOf "https://github.com"
COPY go.mod go.sum /builder/
RUN GOPROXY="direct" go mod tidy
COPY . .
RUN make build

FROM alpine:3.8
WORKDIR /app
COPY --from=base /builder/bin/storage-service /app/storage-service
CMD ["/app/storage-service", "server"]