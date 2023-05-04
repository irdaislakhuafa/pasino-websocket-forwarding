FROM golang:1.18-alpine
WORKDIR /apps
COPY . .
RUN ["go", "run", "./main.go"]