FROM golang:1.12 AS builder
LABEL maintainer="Daniel Lynch <danplynch@gmail.com>"
RUN mkdir -p /go/src/github.com/randomtask1155/hqmon
WORKDIR /go/src/github.com/randomtask1155/hqmon
#ENV GOPATH=/go
#ENV PATH=$GOPATH/bin:$PATH
ADD . .

## need to set CGO_ENABLED=0 probably because of the sendgrid dep
## https://forums.docker.com/t/standard-init-linux-go-195-exec-user-process-caused-no-such-file-or-directory/43777/9
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o hqmon .

FROM scratch
COPY --from=builder /go/src/github.com/randomtask1155/hqmon/hqmon /go/bin/hqmon
ENTRYPOINT ["/go/bin/hqmon"]
