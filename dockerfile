FROM golang:1.21-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o pr-reviewer-service main.go

FROM alpine:latest
WORKDIR /app

COPY --from=build /app/pr-reviewer-service .

ENV PORT=8080

EXPOSE 8080

CMD ["./pr-reviewer-service"]
