FROM e2e-test-base:latest

WORKDIR /portal
ADD docker/sources.tar .
RUN ./mage -v build:installer

WORKDIR /root
COPY .portal.env.yml ./