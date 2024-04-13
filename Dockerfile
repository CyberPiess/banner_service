FROM golang:1.21.5

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /project/go-docker/build/myapp .

EXPOSE 8080

ENTRYPOINT ["/project/go-docker/build/myapp"]
