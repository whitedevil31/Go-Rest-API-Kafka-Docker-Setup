FROM golang:1.16-alpine
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
# ENV MONGO_URI=<>
# ENV PORT=8080
COPY . .
RUN go build -mod=mod producer.go
CMD ["./producer"]