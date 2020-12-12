FROM golang:1.13 

WORKDIR /go/src/app

COPY . .

RUN make core

CMD ["/go/src/app/build/core"]