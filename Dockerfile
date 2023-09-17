FROM scratch
WORKDIR /

ARG USER=1000:1000
USER ${USER}

COPY README.md ./
COPY katastasi /usr/bin/

ENTRYPOINT ["/usr/bin/katastasi"]
CMD []
