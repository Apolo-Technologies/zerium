FROM ubuntu:xenial

RUN \
  apt-get update && apt-get upgrade -q -y && \
  apt-get install -y --no-install-recommends golang git make gcc libc-dev ca-certificates && \
  git clone --depth 1 https://github.com/apolo-technologies/zerium && \
  (cd zerium && make gabt) && \
  cp zerium/build/bin/gabt /gabt && \
  apt-get remove -y golang git make gcc libc-dev && apt autoremove -y && apt-get clean && \
  rm -rf /zerium

EXPOSE 8545
EXPOSE 30303

ENTRYPOINT ["/gabt"]
