FROM golang:1.24.1 AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
