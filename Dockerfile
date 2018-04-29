FROM openlabs/docker-wkhtmltopdf
MAINTAINER Megan O'Keefe <megan037@gmail.com>

RUN mkdir -p /pocketsender/pdf

WORKDIR /pocketsender

COPY bin/pocketsender /bin

ENTRYPOINT ["/bin/pocketsender"]
