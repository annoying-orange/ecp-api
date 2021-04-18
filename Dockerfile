FROM golang:alpine3.13 AS build
RUN apk --no-cache add gcc g++ make git ca-certificates

WORKDIR /go/src/github.com/annoying-orange/ecp
COPY . .
RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 go build -a -v -tags musl -o /go/bin/ecp .

FROM alpine:3.13
WORKDIR /usr/bin
ENV PORT=8080
ENV MYSQL_HOST=localhost:3306
ENV MYSQL_USER=root
ENV MYSQL_PASSWORD=password
COPY --from=build /go/bin .
EXPOSE $PORT
ENTRYPOINT ecp --MYSQL_HOST $MYSQL_HOST --MYSQL_USER $MYSQL_USER --MYSQL_PASSWORD $MYSQL_PASSWORD