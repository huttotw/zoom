FROM golang:1.11
COPY . /go/src/github.com/huttotw/zoom
RUN go install github.com/huttotw/zoom/cmd/zoom
