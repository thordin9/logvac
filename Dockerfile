FROM alpine:3.8
RUN apk --update add tzdata ca-certificates
RUN cp /usr/share/zoneinfo/Europe/Athens /etc/localtime
RUN echo "Europe/Athens" >  /etc/timezone
RUN apk del tzdata
RUN date
RUN apk --update add \
    supervisor \
    git \
    curl \
    unzip \
    nano \
    wget \
    gzip \
    zlib \
    bash
WORKDIR /app
COPY ./dist/logvac /usr/local/bin/
COPY ./supervisord.conf /etc/supervisor/conf.d/supervisord.conf
EXPOSE 6360/tcp
EXPOSE 6361/tcp
EXPOSE 1514/udp
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf","-j","/tmp/supervisord.pid"]
