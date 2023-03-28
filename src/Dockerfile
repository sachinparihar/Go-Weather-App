FROM golang:1.20-alpine

WORKDIR /src

COPY ./src/ .

RUN go mod download

RUN go build -o weather-app .

EXPOSE 8000

CMD ["./weather-app"]
