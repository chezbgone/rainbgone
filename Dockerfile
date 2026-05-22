FROM golang:1.24.1 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/rainbgone .

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/rainbgone /usr/local/bin/rainbgone

EXPOSE 8080

CMD ["/usr/local/bin/rainbgone"]
