FROM golang:1.23.3 AS builder

WORKDIR /usr/local/src

COPY ["app/go.mod", "app/go.sum", "./"]
RUN go mod download

#build
COPY app ./
RUN go build -o ./bin/app cmd/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /

CMD ["/app"]