FROM ubuntu:22.04
RUN apt-get update
WORKDIR /app

COPY ./Klokr /usr/local/bin/klokr

RUN chmod +x /usr/local/bin/klokr


CMD [ "bash" ]

