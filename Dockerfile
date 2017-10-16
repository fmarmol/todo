FROM golang:1.9

WORKDIR /go/src/github.com/fmarmol/todo/
RUN go get  github.com/golang/protobuf/proto\
            github.com/lib/pq\
            golang.org/x/net/context\
            google.golang.org/grpc            
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -i 

FROM golang:1.9-alpine

WORKDIR /root/
COPY --from=0 /go/src/github.com/fmarmol/todo/todo .
EXPOSE 8080
ENTRYPOINT ["./todo"]
