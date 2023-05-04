FROM golang:1.18-alpine
WORKDIR /apps
COPY . .
RUN ["go", "build", "./main.go"]
ENTRYPOINT [ "./main" ]