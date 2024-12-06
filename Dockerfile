ARG GO_VERSION=1

FROM golang:${GO_VERSION}-bookworm as builder
RUN apt-get update && apt-get install -y git
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
ARG GIT_TAG=dev
RUN go build -v -o /run-app .

FROM debian:bookworm
RUN apt-get update && apt-get install -y ca-certificates tzdata && apt-get clean
COPY --from=builder /run-app /usr/local/bin/
EXPOSE 8080
CMD ["run-app"]
