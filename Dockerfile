FROM golang 

ARG REPO

RUN apt-get update -y && apt-get install -y zip

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

ADD . $REPO

WORKDIR $REPO

CMD make package