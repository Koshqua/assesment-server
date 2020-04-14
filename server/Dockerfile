FROM golang:latest
WORKDIR /go/src/server
COPY . . 
RUN go mod download
RUN go install
CMD ["server"]
