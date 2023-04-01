FROM alpine:3

RUN apk add --no-cache gcc git make musl-dev go

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

RUN go install github.com/mattn/go-sqlite3

COPY build/birdbot /usr/bin/birdbot

VOLUME /etc/birdbot

ENTRYPOINT ["/usr/bin/birdbot",  "-c=/etc/birdbot/birdbot.yaml", "-db=/var/lib/birdbot/birdbot.db"]