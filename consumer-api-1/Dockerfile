FROM golang:1.16-alpine
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
# ENV MONGO=some_mobgo
COPY . .
RUN go build -mod=mod consumer.go
CMD ["./consumer"]