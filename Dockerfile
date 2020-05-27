# Stage: build base
FROM golang:1.11.4 AS build_base
WORKDIR /go/src/github.com/thebinary/go-snmp-grpc
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download

# Stage: build
FROM build_base AS build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /server

# Stage: image
FROM alpine:3.6
LABEL maintainer="binary4bytes@gmail.com"
LABEL author="thebinary"
COPY --from=build /server /server
ENTRYPOINT ["/server"]
