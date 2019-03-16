FROM pingcap/alpine-glibc

ADD bin/tidb-server /tidb-server

RUN chmod +x /tidb-server

WORKDIR /

EXPOSE 4000 10080

ENTRYPOINT ["/tidb-server"]
