FROM ubuntu:24.04

WORKDIR /portal

COPY ./portal-installer ./
COPY ./.portal.env.yml ./

ENV PLOG=100
ENV ENABLE_PORTAL_APPHOST_LOG=true
