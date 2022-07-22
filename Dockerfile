# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:latest AS build

WORKDIR /app

COPY . /app

RUN go mod download

RUN go build -o /balance-backend

##
## Deploy
##
FROM gcr.io/distroless/base-debian11:latest

WORKDIR /

COPY --from=build /balance-backend /balance-backend

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/balance-backend"]
