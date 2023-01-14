FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /console

## Deploy
FROM alpine:3.17
WORKDIR /app
COPY --from=build /console console
COPY --from=build /app/templates templates
EXPOSE 80

ENTRYPOINT ["/app/console"]