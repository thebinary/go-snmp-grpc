# Stage: build
FROM golang:1.11.4
ENV GO111MODULE on
WORKDIR /go/src/go-snmp-grpc
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /server

# Stage: image
FROM alpine:3.6
LABEL maintainer="binary4bytes@gmail.com"
LABEL author="thebinary"
COPY --from=0 /server /server
ENTRYPOINT ["/server"]
