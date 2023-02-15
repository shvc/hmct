FROM scratch

COPY hmct /

ENV DEBUG=false \
    MSG='default container message' \
    ADDR=:80

ENTRYPOINT ["/hmct"]
