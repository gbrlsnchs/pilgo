FROM golang:1.14.2-alpine

RUN apk update && apk add git

WORKDIR /src
COPY . .

RUN go install github.com/magefile/mage
CMD mage
