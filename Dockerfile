# Dockerfile
FROM ubuntu:latest
MAINTAINER ish <ish@innogrid.com>

RUN mkdir -p /GraphQL_violin_novnc/
WORKDIR /GraphQL_violin_novnc/

ADD GraphQL_violin_novnc /GraphQL_violin_novnc/
RUN chmod 755 /GraphQL_violin_novnc/GraphQL_violin_novnc

EXPOSE 8001

CMD ["/GraphQL_violin_novnc/GraphQL_violin_novnc"]
