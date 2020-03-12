FROM golang:latest 
RUN mkdir /app 
RUN go get cloud.google.com/go/firestore
ADD . /app/ 
WORKDIR /app 
RUN go build -o mco.fyi . 
CMD ["/app/mco.fyi"]
