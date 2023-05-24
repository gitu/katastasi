FROM alpine:3.17.1
WORKDIR /

COPY README.md ./
COPY katastasi /usr/bin/

ENTRYPOINT ["/usr/bin/katastasi"]
CMD []