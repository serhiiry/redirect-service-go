FROM --platform=linux/x86_64 golang:1.18

WORKDIR /app

COPY app/main.go .
COPY app/redirect-config.json .

ENV GO111MODULE=off

RUN go build -o redirect-service

EXPOSE 80

CMD ["./redirect-service"]
