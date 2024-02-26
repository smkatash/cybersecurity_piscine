FROM debian:bookworm

RUN apt-get update && apt-get -y install golang \
    build-essential \
    libpcap-dev \
    iproute2 \
    ca-certificates  \
    iputils-ping

WORKDIR /inquisitor

CMD [ "tail", "-f", "/dev/null" ]