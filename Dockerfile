FROM alpine:3.4

RUN apk update && \
  apk add \
    ca-certificates && \
  rm -rf /var/cache/apk/*
EXPOSE 8080
ADD rickjames /bin/
ENTRYPOINT ["/bin/rickjames"]
