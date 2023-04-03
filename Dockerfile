FROM golang:1.20-bullseye

RUN apt update -y

COPY build/birdbot /usr/bin/birdbot

VOLUME /etc/birdbot

ENTRYPOINT ["/usr/bin/birdbot",  "-c=/etc/birdbot/birdbot.yaml", "-db=/var/lib/birdbot/birdbot.db"]