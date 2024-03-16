FROM golang:latest

WORKDIR /app

COPY go.mod .

COPY . .

RUN go mod tidy

RUN go build -o -v api
 
EXPOSE 8000

CMD ["./api"]