FROM golang:1.16

WORKDIR /urlshortener

COPY src/. ./

RUN go mod download
RUN go get ./storage/

RUN go build -o /urlshortener-app

CMD [ "/urlshortener-app" ]