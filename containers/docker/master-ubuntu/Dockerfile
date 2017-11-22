FROM ubuntu:xenial

RUN \
  apt-get update && apt-get upgrade -q -y && \
  apt-get install -y --no-install-recommends golang git make gcc libc-dev ca-certificates && \
  git clone --depth 1 --branch release/1.7 https://github.com/abt/zerium && \
  (cd zerium && make gzrm) && \
  cp zerium/build/bin/gzrm /gzrm && \
  apt-get remove -y golang git make gcc libc-dev && apt autoremove -y && apt-get clean && \
  rm -rf /zerium

EXPOSE 8545
EXPOSE 32310

ENTRYPOINT ["/gzrm"]
