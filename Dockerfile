FROM golang:1.16
LABEL maintainer="forum project"
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go get github.com/mattn/go-sqlite3
RUN go get golang.org/x/crypto/bcrypt
RUN go get github.com/satori/go.uuid
RUN go build -o main .
CMD ["./main"]