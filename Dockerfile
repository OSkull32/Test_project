FROM golang:latest as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build  -o main .

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]