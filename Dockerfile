# Pull base image
FROM resin/rpi-raspbian:jessie

VOLUME ["/config"]

# Install OpenJDK 8 and dependencies
COPY jdk-8u121-fcs-bin-b13-linux-arm32-vfp-hflt-12_dec_2016.tar.gz /

RUN tar zxvf /jdk-8u121-fcs-bin-b13-linux-arm32-vfp-hflt-12_dec_2016.tar.gz -C /opt && \
	update-alternatives --install "/usr/bin/java" "java" "/opt/jdk1.8.0_121/bin/java" 1 && \
	update-alternatives --install "/usr/bin/javac" "javac" "/opt/jdk1.8.0_121/bin/javac" 1 

# Set JAVA_HOME variable
ENV JAVA_HOME /opt/jdk1.8.0/

RUN apt-get update && \
	apt-get install -y git wget

RUN apt-get install -y curl

COPY firstrun.sh /

RUN chmod +x /firstrun.sh
