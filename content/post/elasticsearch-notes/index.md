---
title: "Elasticsearch Notes"
description: 
date: 2021-07-19T18:22:26+08:00
image: 
math: 
license: 
hidden: false
draft: true
category:
    -: 数据库
tag:
    -:  Elasticsearch
---

本文基于Elasticsearch v7.13.3版本

## install

利用docker快速部署一个本地环境

### docker

安装ES

```shell
docker network create elastic
docker pull docker.elastic.co/elasticsearch/elasticsearch:7.13.3
docker run --name es01-test --net elastic -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:7.13.3
```

安装kibana

```shell
docker pull docker.elastic.co/kibana/kibana:7.13.3
docker run --name kib01-test --net elastic -p 5601:5601 -e "ELASTICSEARCH_HOSTS=http://es01-test:9200" docker.elastic.co/kibana/kibana:7.13.3
```



## CRUD

### Add Data





