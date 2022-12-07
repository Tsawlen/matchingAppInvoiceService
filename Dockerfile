# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /matchingAppInvoiceService

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /invoice-service

EXPOSE 8084

CMD [ "/invoice-service" ]