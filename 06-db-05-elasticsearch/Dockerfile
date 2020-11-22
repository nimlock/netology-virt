FROM centos:7

RUN yum install -y wget

ENV ELASTIC_VER=7.10.0 \
    ES_CUSTOM_USER=elasticuser \
    ES_PATH_DATA=/var/lib/elasticsearch-data

RUN adduser -U ${ES_CUSTOM_USER}

WORKDIR /app

RUN wget -q https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-${ELASTIC_VER}-linux-x86_64.tar.gz && \
    wget -q https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-${ELASTIC_VER}-linux-x86_64.tar.gz.sha512 && \
    sha512sum -c elasticsearch-${ELASTIC_VER}-linux-x86_64.tar.gz.sha512 && \
    tar -xzf elasticsearch-${ELASTIC_VER}-linux-x86_64.tar.gz && \
    rm elasticsearch-${ELASTIC_VER}-linux-x86_64.tar.gz* && \
    chown -R ${ES_CUSTOM_USER} .

RUN mkdir -p ${ES_PATH_DATA} && chown ${ES_CUSTOM_USER} ${ES_PATH_DATA}

USER ${ES_CUSTOM_USER}

WORKDIR /app/elasticsearch-${ELASTIC_VER}/

EXPOSE 9200/tcp

CMD ./bin/elasticsearch
