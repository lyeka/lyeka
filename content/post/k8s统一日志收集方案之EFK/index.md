---
title: "K8S日志管理方案之EFK"
description: 
date: 2021-07-06T20:18:45+08:00
image: 
math: 
license: 
hidden: false
draft: true
---

## 常见日志管理方案

在统一日志管理中，比较出名的是ELK方案，即 Elasticsearch、Logstash 和 Kibana组合。Logstash 负责日志的收集，可以从多方数据源收集将其上报至Elasticsearch；Elasticsearch是一个搜索和分析引擎，通过ES我们可以全文搜索日志以快速定位问题；Kibana则是ES的可视化界面，通过Kibana可以对数据绘制报表图。

在K8S中，比较推荐的是EFK方案，即 Elasticsearch、Fluentd 和 Kibana组合。相对于ELK，只是日志收集组件变成了Fluentd 。

所以ELK vs EFK 可以看成 Logstash  vs Fluentd 。但除了Fluentd， CNCF还推出另外一款日志收集组件Fluent Bit，Fluent Bit其实是Fluentd的子项目，相对于Fluentd更加轻量，不过因为年限不久等原因，功能上相对Fluentd少了一些。

本文主要介绍以**Fluent Bit为日志收集组件的EFK方案**



###  Logstash  vs Fluentd 

TODO

### Fluentd vs Fluentbit

|          | Fluentd       | Fluent Bit                |
| -------- | ------------- | ------------------------- |
| 开发语言 | C & Ruby      | C                         |
| 依赖     | 依赖 Ruby Gem | 无依赖，除非一些插件需要  |
| 内存占用 | 大约40MB      | 大约650KB                 |
| 插件     | 1000+插件     | 70+插件                   |
| 使用范围 | 容器、服务器  | 容器、服务器、嵌入式Linux |
|          |               |                           |

## K8S日志收集架构

K8S本身不提供原生的日志管理方案，日志无论是输出到stdout/stderr还是写入文件，如果没有作持久化处理的话，当Pod被驱逐，Node宕机等情况下，日志都会丢失。所以需要开发人员提供一种日志收集方案以使日志的生命周期与Pod、Node相互独立。

K8S 常见的日志收集架构有：

- 工作节点日志代理
- sidecar 容器日志代理
- 应用直接暴露日志

### 工作节点日志代理

![Using a node level logging agent](https://d33wubrfki0l68.cloudfront.net/2585cf9757d316b9030cf36d6a4e6b8ea7eedf5a/1509f/images/docs/user-guide/logging/logging-with-node-agent.png)

通过 DaemonSet 将 logging-agent (日志收集组件)以pod的形式部署到每一个节点上，logging-agent 会监控日志文件（K8S会将自动节点上所有 Pod 的 stdout 和 stder r重定向该日志文件），将日志发送到 Logging Backend（如ES，Fluentd，S3等）

### sidecar容器日志代理

![带数据流容器的边车容器](https://d33wubrfki0l68.cloudfront.net/5bde4953b3b232c97a744496aa92e3bbfadda9ce/39767/images/docs/user-guide/logging/logging-with-streaming-sidecar.png)

通过在 pod 里面部署多一个容器（称为 sidecar ）运行 logging-agent 



### 应用直接暴露日志



ref

- [Fluent Bit官方文档](https://docs.fluentbit.io/manual/)
- [K8S日志架构](https://kubernetes.io/docs/concepts/cluster-administration/logging/)

