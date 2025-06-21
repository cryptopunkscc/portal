FROM ubuntu:24.04

WORKDIR /portal

COPY ./portal-installer ./

RUN ./portal-installer

ENV PLOG=100

CMD ["sh", "-c", "portal"]
