FROM golang:alpine AS builder

WORKDIR /app

# copy and install dependencies
COPY go.mod go.sum ./
RUN GOOS=linux go mod download

COPY . .

RUN GOOS=linux go build -o main .

FROM alpine

COPY --from=builder /app/main /application

RUN chmod +x /application

ENTRYPOINT /application