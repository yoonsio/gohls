FROM golang:1.20.2 as builder

workdir /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gohls .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/gohls /usr/local/bin/gohls

ENTRYPOINT ["/usr/local/bin/gohls"]

