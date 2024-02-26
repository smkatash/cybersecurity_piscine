FROM debian:bookworm

RUN apt-get update && apt-get -y install golang \
    build-essential 

# TODO for eval add some dummy folder
RUN mkdir -p infection 


CMD [ "tail", "-f", "/dev/null" ]