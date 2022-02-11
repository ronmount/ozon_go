FROM golang:1.17.6-alpine3.15 as builder
COPY go.mod go.sum /go/src/github.com/ronmount/ozon_go/
WORKDIR /go/src/github.com/ronmount/ozon_go
RUN go mod download
COPY . /go/src/github.com/ronmount/ozon_go
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o build/ozon_go ./cmd/main.go

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/ronmount/ozon_go/build/ozon_go /usr/bin/ozon_go
COPY --from=builder /go/src/github.com/ronmount/ozon_go/.env .
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/ozon_go"]