FROM golang:1.25.4

ARG MUT_VERSION=0.3.0

RUN mkdir -p /app/configs
RUN mkdir -p /app/var/logs
RUN apt-get update

WORKDIR /app

RUN curl -sL https://github.com/Clivern/Mut/releases/download/v${MUT_VERSION}/mut_Linux_x86_64.tar.gz | tar xz
RUN rm LICENSE
RUN rm README.md

COPY ./config.dist.yml /app/configs/

EXPOSE 8080

VOLUME /app/configs
VOLUME /app/var

RUN ./mut version

CMD ["./mut", "server", "-c", "/app/configs/config.dist.yml"]
