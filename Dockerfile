FROM golang:1.13

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUn go build
#RUN go install -v ./...

CMD ["./shorty"]
