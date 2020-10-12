FROM golang:alpine
ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 CI=true
WORKDIR /app
COPY . .
ENTRYPOINT ["go", "test", "-v", "./app/test", "-coverprofile", "cover.out"]