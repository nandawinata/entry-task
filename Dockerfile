FROM golang:latest
RUN go get -u github.com/golang/dep/cmd/dep \
&& mkdir /go/src/github.com/nandawinata \
&& git clone https://github.com/nandawinata/entry-task /go/src/github.com/nandawinata/entry-task

WORKDIR /go/src/github.com/nandawinata/entry-task/

RUN dep ensure -v
RUN go build -o /go/src/github.com/nandawinata/entry-task/cmd/app

ENTRYPOINT ["/go/src/github.com/nandawinata/entry-task/"]