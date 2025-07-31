ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY litefs.yml .

COPY . .
RUN go build -v -o /run-app .


FROM debian:bookworm

COPY --from=builder /run-app /usr/local/bin/
# CMD ["run-app"]

### LiteFS
RUN apt-get update -y && apt-get install -y ca-certificates fuse3 sqlite3 curl

COPY --from=flyio/litefs:0.5 /usr/local/bin/litefs /usr/local/bin/litefs

### Copy LiteFS cfg
COPY --from=builder /usr/src/app/litefs.yml /etc/litefs.yml

ENTRYPOINT litefs mount

# Install dotenvx
# RUN apt-get install curl
RUN curl -fsS https://dotenvx.sh/install.sh | sh
# COPY --from=builder /app/.env* ./

# COPY --from=builder /app/main .
# EXPOSE 1323

# Prepend dotenvx run
CMD ["dotenvx", "run", "--", "run-app"]
