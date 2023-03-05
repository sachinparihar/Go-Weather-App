FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN go mod init main

RUN go mod download

RUN go build -o weather-app .

EXPOSE 8080

CMD ["./weather-app"]
