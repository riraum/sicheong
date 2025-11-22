ARG GO_VERSION=1
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o /run-app .

FROM --platform=linux/amd64 debian:bookworm
COPY --from=builder /run-app /usr/local/bin/
CMD ["run-app"]


# To remove or adjust for GKE
### LiteFS
# RUN apt-get update -y && apt-get install -y ca-certificates fuse3 sqlite3

# COPY --from=flyio/litefs:0.5 /usr/local/bin/litefs /usr/local/bin/litefs

# ### Copy LiteFS cfg
# COPY --from=builder /usr/src/app/litefs.yml /etc/litefs.yml

# ENTRYPOINT litefs mount
