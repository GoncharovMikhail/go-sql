FROM golang:latest as BUILD
MAINTAINER https://github.com/GoncharovMikhail
USER GoncharovMikhail
RUN apt-get clean && \
    apt-get update && \
    apt-get install -y binutils upx
WORKDIR builddir
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o go-gmux-db-testcontainers

FROM scratch
COPY --from=BUILD ./go/builddir/go-gmux-db-testcontainers .
ENTRYPOINT ["./go-gmux-db-testcontainers"]