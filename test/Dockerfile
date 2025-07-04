FROM ubuntu:24.04

WORKDIR /portal

COPY ./portal-installer ./
COPY ./.portal.env.yml ./

RUN ./portal-installer

ENV PLOG=100
ENV ENABLE_PORTAL_APPHOST_LOG=true
