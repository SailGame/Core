FROM golang:1.18

WORKDIR /go/src/app

COPY build/core .

CMD ["/go/src/app/core"]