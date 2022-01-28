FROM golang:1.17-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o cat-app ./main.go

CMD ["./cat-app"]