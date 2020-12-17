FROM cmulk/wireguard-docker:alpine

#Installation of openresolv fix the DNS issue with wireguard and docker
#https://github.com/linuxserver/docker-wireguard/issues/15
RUN apk update && apk add bind-tools curl openresolv

COPY run.sh /scripts
COPY ./scripts/wireguard-connection-validation.sh /scripts

ENV IPTABLES_MASQ=1

WORKDIR /scripts

CMD ["run.sh"]
