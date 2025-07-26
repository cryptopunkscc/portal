FROM e2e-test-base:latest

WORKDIR /portal
RUN ./mage -v build:installer

WORKDIR /root
COPY .portal.env.yml ./