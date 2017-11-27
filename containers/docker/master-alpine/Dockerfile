FROM alpine:3.5

RUN \
  apk add --update go git make gcc musl-dev linux-headers ca-certificates && \
  git clone --depth 1 --branch release/1.7 https://github.com/apolo-technologies/zerium && \
  (cd zerium && make gzrm) && \
  cp zerium/build/bin/gzrm /gzrm && \
  apk del go git make gcc musl-dev linux-headers && \
  rm -rf /zerium && rm -rf /var/cache/apk/*

EXPOSE 2008
EXPOSE 32310

ENTRYPOINT ["/gzrm"]
