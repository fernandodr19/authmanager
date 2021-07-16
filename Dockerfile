FROM golang:1.16.3-alpine3.13

WORKDIR /app

COPY Makefile go.mod go.sum ./

RUN go mod download

COPY . .

RUN make compile

EXPOSE 3000

CMD ["./build/app"]