FROM golang:latest 
RUN mkdir /app 
RUN go get cloud.google.com/go/firestore
ADD . /app/ 
WORKDIR /app 
RUN go build -o redir . 
CMD ["/app/redir"]
