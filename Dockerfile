FROM golang:lateset

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /BaiDuPanLoad

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

EXPOSE 8080

RUN go build -o panload .

ENTRYPOINT ["./panload"]

