FROM golang 

ENV GOOS=linux GO111MODULE=on

ARG REPO

RUN apt-get update -y && apt-get install -y zip

ADD . $REPO

WORKDIR $REPO

CMD make package