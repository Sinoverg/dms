FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

# Build your Go binary for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o /dbms

EXPOSE 1488

CMD ["/dbms"]