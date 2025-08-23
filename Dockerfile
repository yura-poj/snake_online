FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o snake-online ./main_package

EXPOSE 80

CMD ["./snake-online"]