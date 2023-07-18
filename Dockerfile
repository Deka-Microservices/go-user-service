FROM golang:1.20 AS builder

WORKDIR /build

ARG VERSION=dev

COPY . .
RUN go mod download
RUN go build -ldflags="-X 'main.Version=${VERSION}'" -o ./server cmd/server/main.go 

FROM ubuntu:latest
WORKDIR /app 
COPY --from=builder /build/server /app/

EXPOSE 8080

CMD [ "/app/server" ]


