FROM alpine:3.3

ADD *.go /subjects-transformer/
ADD handlers/*.go /subjects-transformer/handlers/
ADD service/* /subjects-transformer/service/
ADD model/*.go /subjects-transformer/model/

RUN apk add --update bash \
  && apk --update add git bzr \
  && apk --update add go \
  && export GOPATH=/gopath \
  && REPO_PATH="github.com/Financial-Times/subjects-transformer" \
  && mkdir -p $GOPATH/src/${REPO_PATH} \
  && cp -r subjects-transformer/* $GOPATH/src/${REPO_PATH} \
  && cd $GOPATH/src/${REPO_PATH} \
  && go get -t ./... \
  && go build \
  && mv subjects-transformer /app \
  && apk del go git bzr \
  && rm -rf $GOPATH /var/cache/apk/*
CMD [ "/app" ]