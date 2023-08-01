FROM scratch
WORKDIR /

COPY README.md ./
COPY katastasi /usr/bin/

ENTRYPOINT ["/usr/bin/katastasi"]
CMD []