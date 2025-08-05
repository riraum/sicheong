ARG GO_VERSION=1

### Builder stage ###
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy LiteFS config and encrypted env files explicitly
COPY litefs.yml .
COPY .env.production .

# Copy all source files and build the app binary
COPY . .
RUN go build -v -o /run-app .

### Runtime stage ###
FROM debian:bookworm

# Install runtime dependencies first
RUN apt-get update -y && apt-get install -y ca-certificates fuse3 sqlite3 curl

# Install dotenvx
RUN curl -fsS https://dotenvx.sh/install.sh | sh

# Copy LiteFS binary
COPY --from=flyio/litefs:0.5 /usr/local/bin/litefs /usr/local/bin/litefs

# Copy the built Go binary and make executable
COPY --from=builder /run-app /usr/local/bin/run-app
RUN chmod +x /usr/local/bin/run-app

# Copy LiteFS config file
COPY --from=builder /usr/src/app/litefs.yml /etc/litefs.yml

# Set working directory for your app
WORKDIR /usr/src/app

# Copy encrypted env files
COPY --from=builder /usr/src/app/.env* /usr/src/app/

ENTRYPOINT ["sh", "-c", "litefs mount & sleep 1 && exec \"$@\"", "--"]
CMD ["dotenvx", "run", "-f", "/usr/src/app/.env.production", "--", "run-app"]
