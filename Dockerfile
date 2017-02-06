# Pull base image
FROM resin/rpi-raspbian:jessie
MAINTAINER Talmai Oliveira <to@talm.ai>

VOLUME ["/config"]

# Install OpenJDK 8 runtime without X11 support
RUN echo "deb http://ftp.debian.org/debian jessie-backports main" | sudo tee /etc/apt/sources.list.d/backports.list && \
    apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 8B48AD6246925553 && \
    apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 7638D0442B90D010 && \
    apt-get update && \
    apt-get -t jessie-backports install -y openjdk-8-jre-headless --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

# Set JAVA_HOME variable
ENV JAVA_HOME /usr/lib/jvm/java-8-openjdk-armhf

COPY firstrun.sh /etc/my_init.d/firstrun.sh

RUN chmod +x /etc/my_init.d/firstrun.sh