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
FROM bash:5.2.15-alpine3.17
COPY --from=build /go/bin/aproxy /
CMD ["/aproxy"]



#docker build . -t ghcr.io/katasec/aproxy:v0.0.1
#docker run -e APROXY_TARGET_URL="https://go.dev" -it ghcr.io/katasec/aproxy:v0.0.1