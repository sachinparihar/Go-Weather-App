FROM golang:1.20-alpine

WORKDIR /app

# Set the environment variable for the Docker environment
ENV FRONTEND_DIR=/app/src/Frontend/

COPY . .

RUN cd src/Backend && go build -o ../../weather-app .

EXPOSE 8000

CMD ["./weather-app"]