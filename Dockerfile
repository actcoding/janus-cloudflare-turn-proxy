FROM golang:1.23-alpine AS build

RUN apk add --no-cache \
        curl \
        git \
    && \
    sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -b /usr/local/bin

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN task build:release


FROM alpine:3

ENV CONTAINER=true

COPY --from=build /usr/src/app/dist/jctp /app

HEALTHCHECK CMD sh -c "[ ! -f /tmp/failure ]"

USER 1001

CMD ["/app"]
