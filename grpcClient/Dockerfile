FROM golang:latest

WORKDIR /
COPY . .

RUN go mod download

EXPOSE 3001

CMD ["go", "run", "Main.go"]