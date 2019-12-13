FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

RUN go get -u github.com/golang/dep/cmd/dep \
&& mkdir /go/src/github.com/nandawinata \
&& git clone https://github.com/nandawinata/entry-task /go/src/github.com/nandawinata/entry-task

WORKDIR /go/src/github.com/nandawinata/entry-task/
RUN dep ensure -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main github.com/nandawinata/entry-task/cmd/app

FROM mysql:5.7
WORKDIR /root/
COPY --from=builder /go/src/github.com/nandawinata/entry-task/scripts/sql ./sql
RUN chmod a+x sql/init_table.sh
RUN ls -a
RUN cat sql/init_table.sql
RUN sql/init_table.sh

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/nandawinata/entry-task/main .
COPY --from=builder /go/src/github.com/nandawinata/entry-task/configs ./configs
COPY --from=builder /go/src/github.com/nandawinata/entry-task/.env .
EXPOSE 8080
CMD ["./main"]