FROM ubuntu:24.04

RUN apt-get update && \
    apt-get install -y nodejs npm && \
    apt-get install -y libwebkit2gtk-4.1-0 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN npm -v
