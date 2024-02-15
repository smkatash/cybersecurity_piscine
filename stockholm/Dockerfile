FROM debian:bookworm

RUN apt-get update && apt-get install sudo

RUN sudo apt install -y tar gzip build-essential wget
RUN sudo wget https://golang.org/dl/go1.22.0.linux-amd64.tar.gz && sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz \
    && sudo rm go1.22.0.linux-amd64.tar.gz

ENV PATH="$PATH:/usr/local/go/bin"
ENV GOPATH="$HOME/go"
ENV PATH="$PATH:$GOPATH/bin"
RUN mkdir -p infection


CMD [ "tail", "-f", "/dev/null" ]