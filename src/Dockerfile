FROM golang:1.20-alpine

WORKDIR /src

COPY . .

COPY main /app

RUN go mod init main

RUN go mod download

RUN go build -o weather-app .

EXPOSE 8000

CMD ["./weather-app"]
