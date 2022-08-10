FROM golang:1.17-alpine AS BUILD
WORKDIR /app
COPY . .
RUN go build -o app main.go

FROM alpine:3.14
WORKDIR /app
COPY --from=BUILD /app/app .
COPY config.yaml .
EXPOSE 8000

CMD ["./app"]

