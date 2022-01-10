# syntax=docker/dockerfile:1

From golang:1.17.5-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o index .

EXPOSE 8001

CMD [ "./index" ]
