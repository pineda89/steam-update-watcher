FROM golang

RUN go get -u gopkg.in/telegram-bot-api.v4
RUN go get -u github.com/sirupsen/logrus

ADD . /go/src/steam-update-watcher

RUN go install steam-update-watcher

ENTRYPOINT /go/bin/steam-update-watcher