FROM docker.io/library/golang:1.23.7

RUN apt-get update && \
    apt-get install -y nodejs npm && \
#    apt-get install -y libwebkit2gtk-4.1-0 && \
    apt-get install -y libwebkit2gtk-4.1-dev && \
    apt-get install -y gcc libgtk-3-dev libayatana-appindicator3-dev && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

ENV PLOG=100
ENV ENABLE_PORTAL_APPHOST_LOG=true