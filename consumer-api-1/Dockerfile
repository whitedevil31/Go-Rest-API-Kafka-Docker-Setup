FROM golang:1.16-alpine
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
ENV MONGO=
ENV MAILER_DOMAIN=
ENV MAILER_KEY=
ENV MAILER_ID=
ENV SPREADSHEET_ID=
ENV SPREADSHEET_RANGE=
COPY . .
RUN go build -mod=mod consumer.go
CMD ["./consumer"]