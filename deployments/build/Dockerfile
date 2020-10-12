FROM golang:alpine AS stage1
ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 
WORKDIR /app
COPY . .
RUN go build -o api ./cmd/main.go

FROM alpine
WORKDIR /app
RUN apk add mysql-client
COPY --from=stage1 /app/api .
COPY --from=stage1 /app/migrations/build ./migrations
EXPOSE 8081
CMD ["./api"]  
