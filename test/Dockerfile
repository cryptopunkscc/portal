FROM e2e-test-base:latest

WORKDIR /portal

COPY ./install-portal-to-astral ./
COPY ./.portal.env.yml ./

ENV PLOG=100
ENV ENABLE_PORTAL_APPHOST_LOG=true
