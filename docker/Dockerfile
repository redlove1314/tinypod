FROM hehety/golang:godfs
RUN mkdir /tinypod && cd /tmp && rm -rf tinypod && git clone https://github.com/hetianyi/tinypod.git && \
    cd tinypod && \
    ./make.sh && cp bin/* /tinypod && \
    cd / && tar czvf tinypod.tar.gz tinypod


FROM alpine:latest
COPY --from=0 /tinypod.tar.gz /
RUN tar -xzf /tinypod.tar.gz -C / && rm -f /tinypod.tar.gz && \
    ln -s /tinypod/http /usr/bin/http && \
    ln -s /tinypod/proxy /usr/bin/proxy
WORKDIR /tinypod
