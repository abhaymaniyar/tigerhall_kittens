FROM golang:1.18 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o tigerhall_kittens

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
FROM debian:buster-slim
RUN set -eux; apt-get update; apt-get install -y --no-install-recommends \
    ca-certificates \
    netbase \
    ;

COPY --from=builder /app/tigerhall_kittens /tigerhall_kittens
CMD ["/tigerhall_kittens"]