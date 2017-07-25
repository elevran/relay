FROM alpine:3.6

MAINTAINER Etai Lev Ran <elevran@google.com>

WORKDIR /opt/relayd/

RUN apk update && apk upgrade \
  && apk add ca-certificates \
  && rm -rf /var/cache/apk/*

ADD bin/relayd /opt/relayd/

USER nobody:nobody
ENTRYPOINT ["/opt/relayd/relayd"]
