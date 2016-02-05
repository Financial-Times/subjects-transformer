FROM alpine

ENV GOPATH /gopath
ENV REPO_PATH "github.com/Financial-Times/subjects-transformer"

RUN apk add --update bash \
  && apk --update add git bzr \
  && echo "http://dl-4.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories \
  && apk --update add go \
  && mkdir -p $GOPATH/src/${REPO_PATH}

ADD . $GOPATH/src/$REPO_PATH/

RUN cd $GOPATH/src/${REPO_PATH} \
  && go get -t ./... \
  && go test ./... \
  && go build \
  && mv subjects-transformer /subjects-transformer \
  && apk del go git bzr \
  && rm -rf $GOPATH /var/cache/apk/*

CMD /subjects-transformer -baseUrl=$BASE_URL -structureServiceBaseUrl=$STRUCTURE_SERVICE_BASE_URL -structureServiceUsername=$USER -structureServicePassword=$PASS -structureServicePrincipalHeader=$PRINCIPAL_HEADER
