FROM debian:bookworm

RUN apt-get update && apt-get install dpkg --add-architecture i386 && apt-get update
# install 32-bit libraries
RUN apt-get install -y libc6:i386 \
                            gdb \
                            build-essential
RUN apt-get install -y gcc-multilib 
# compile gcc -m32 source.c -o source

CMD [ "tail", "-f", "/dev/null" ]