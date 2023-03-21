FROM golang:alpine AS build

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o acceptance-cmd ./tests/acceptance

FROM alpine:latest

COPY --from=build /app/acceptance-cmd /usr/local/bin/

CMD ["acceptance-cmd"]