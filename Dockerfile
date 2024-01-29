# docker build "$PWD" --build-arg service="customer" --build-arg commit="$(git rev-parse --short HEAD)" --build-arg version=v1.0.0 -t microservices-poc/customer:1.0.0
# docker build "$PWD" --build-arg service="customer" --build-arg commit="$(git rev-parse --short HEAD)" --build-arg version="$(git rev-parse --short HEAD)" -t microservices-poc/customer-pre-release:"$(git rev-parse --short HEAD)"

# Stage 1: modules caching
FROM golang:1.21-alpine as modules
LABEL maintainer="htquangg@gmail.com"

WORKDIR /microservices-poc

COPY go.* .

RUN go mod download

# Stage 2: build
FROM golang:1.21-alpine as builder
LABEL maintainer="htquangg@gmail.com"

ARG service
ARG version
ARG commit

ENV GOOSE linux
ENV GOARCH amd64
ENV CGO_ENABLED 0
ENV GO111MODULE=on

WORKDIR /microservices-poc

COPY --from=modules /go/pkg /go/pkg
COPY . .

RUN go build -o ./$service -trimpath -ldflags "-s -w -X main.version=$version -X main.commitID=$commit" ./internal/services/$service \
    && cp ./internal/services/$service/config/config.development.yaml ./config.yaml

# Stage 3: deploy
FROM alpine:3 as runtime
LABEL maintainer="htquangg@gmail.com"

ARG service
ARG version

LABEL service=$service
LABEL version=$version

WORKDIR /microservices-poc

RUN apk update \
    && apk --no-cache add \
        bash \
    && echo "UTC" > /etc/timezone

ENV TZ UTC
ENV CONFIG_PATH /microservices-poc/config.yaml

COPY --from=builder /microservices-poc/config.yaml ./config.yaml
COPY --from=builder /microservices-poc/$service ./app

EXPOSE 30001

CMD ["/microservices-poc/app"]
