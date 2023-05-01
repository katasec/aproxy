# ------------------------------------------------------------------------------------------
# Build App in Golang Container
# ------------------------------------------------------------------------------------------

FROM golang:1.20.3-alpine as build
WORKDIR /go/src/app
COPY . .

RUN go mod download && \
    CGO_ENABLED=0 go build -o /go/bin/aproxy


# ------------------------------------------------------------------------------------------
# Copy compiled binary on to golang distro
# ------------------------------------------------------------------------------------------
# FROM bash:5.2.15-alpine3.17
FROM ghcr.io/katasec/tailscale:0.07
COPY --from=build /go/bin/aproxy /
COPY tailscale.sh /usr/local/bin/tailscale.sh
# CMD ["/aproxy"]
ENTRYPOINT ["docker-entrypoint.sh"]


#docker build . -t ghcr.io/katasec/aproxy:v0.0.2
#docker run -it ghcr.io/katasec/aproxy:v0.0.2
#docker run -e APROXY_TARGET_URL="https://go.dev" -e APROXY_TARGET_PORT="1337" -it ghcr.io/katasec/aproxy:v0.0.2