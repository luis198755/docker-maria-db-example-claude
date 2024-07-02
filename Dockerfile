# Dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app

#COPY go.mod go.sum ./
#RUN go mod download

#COPY *.go ./
# Copy the source from the current directory to the Working Directory inside the container
COPY ./main.go ./main.go

# Download any dependencies
RUN go mod init ecommerce-api

# Add the required dependencies
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /ecommerce-api

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /ecommerce-api .

EXPOSE 8080

CMD ["./ecommerce-api"]