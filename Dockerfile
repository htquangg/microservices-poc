# docker build "$PWD" --build-arg service="customer" --build-arg commit="$(git rev-parse --short HEAD)" --build-arg version=v1.0.0 -t microservices-poc/customer:1.0.0
# docker build "$PWD" --build-arg service="customer" --build-arg commit="$(git rev-parse --short HEAD)" --build-arg version="$(git rev-parse --short HEAD)" -t microservices-poc/customer-pre-release:"$(git rev-parse --short HEAD)"

# Stage 1: build
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

COPY . .

RUN go build -o ./$service -trimpath -ldflags "-s -w -X main.version=$version -X main.commitID=$commit" ./internal/services/$service
RUN cp ./internal/services/$service/config/config.docker.yaml ./config.$service.yaml

# Stage 2: deploy
FROM alpine:3 as runtime
LABEL maintainer="htquangg@gmail.com"

ARG service
ARG commit
ARG version

LABEL service=$service
LABEL version=$version

WORKDIR /microservices-poc

RUN apk update \
    && apk --no-cache add \
        bash \
    && echo "UTC" > /etc/timezone

ENV TZ UTC
ENV CONFIG_PATH /microservices-poc/config.$service.yaml

COPY --from=builder /microservices-poc/scripts/entrypoint.sh ./entrypoint.sh
COPY --from=builder /microservices-poc/config.$service.yaml ./config.$service.yaml
COPY --from=builder /microservices-poc/$service ./$service

RUN chmod 755 ./entrypoint.sh

ENV SERVICE_NAME $service

ENTRYPOINT ["/microservices-poc/entrypoint.sh"]
