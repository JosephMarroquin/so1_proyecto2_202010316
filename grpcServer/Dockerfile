FROM golang:latest

WORKDIR /
COPY . .

RUN go mod download

EXPOSE 3000

CMD ["go", "run", "Main.go"]