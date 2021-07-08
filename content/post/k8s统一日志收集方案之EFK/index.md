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

本文主要介绍以**Fluent Bit为日志收集组件的EFK方案**。



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

#### 用 Fluentd 还是 Fluent Bit ？

无论是 Fluentd 还是 Fluent Bit 都可以单独作为日志收集组件工作。如果你需要的插件功能只有 Fluentd 才有，那只能选择 Fluentd ，如果有内存节省的需求可以考虑使用 Fluentd + Fluent Bit 组合工作， Fluent Bit 以 DaemonSet 部署到各个工作节点上负责日志收集，转发到 Fluentd 作日志的聚合以及处理

。如下图示

![img](https://dytvr9ot2sszz.cloudfront.net/wp-content/uploads/2018/06/kuberbetes-monitoring-arch-1.jpg)



最初我们采用的就是 Fluentd + Fluent Bit 的组合，不过后面分析了我们的使用场景 Fluent Bit 已经满足，考虑到**最小化依赖的原则以及内存的节省**，就移除了 Fluentd ，Fluent Bit直接作日志的收集、聚合、处理，后上传到 ES 和 S3。



## K8S日志收集架构

K8S 本身不提供原生的日志管理方案，日志无论是输出到 stdout/stderr 还是写入文件，如果没有作持久化处理的话，当 pod 被驱逐，node 宕机等情况下，日志都会丢失。所以需要开发人员提供一种日志收集方案以使日志的生命周期与 pod 、node 相互独立。 

K8S 常见的日志收集架构有：

- 工作节点日志代理
- sidecar 容器日志代理
- 应用直接暴露日志

### 工作节点日志代理

![Using a node level logging agent](https://d33wubrfki0l68.cloudfront.net/2585cf9757d316b9030cf36d6a4e6b8ea7eedf5a/1509f/images/docs/user-guide/logging/logging-with-node-agent.png)

通过 DaemonSet 将 logging-agent (日志收集组件)以pod的形式部署到每一个节点上，logging-agent 会监控日志文件（K8S会将自动节点上所有 Pod 的 stdout 和 stderr 重定向该日志文件），将日志发送到 Logging Backend（如 ES，Fluentd，S3 等）。

EFK方案即这种架构。

### sidecar容器日志代理

![带数据流容器的边车容器](https://d33wubrfki0l68.cloudfront.net/5bde4953b3b232c97a744496aa92e3bbfadda9ce/39767/images/docs/user-guide/logging/logging-with-streaming-sidecar.png)

sidecar 容器日志代理同样需要一个节点级别的日志代理工作，这种场景主要是处理应用无法输出到 stdout/stderr的情形。

通过在 pod 里面部署多一个容器（称为 sidecar ）运行 logging-agent ，logging-agent 将应用输出的日志（可以是文件、套接字，journald 等）输出到 sidecar 的 stdout/stderr ，节点的日志代理会捕获到这些日志将其发送至 Logging Backend。



### 应用直接暴露日志

![Exposing logs directly from the application](https://d33wubrfki0l68.cloudfront.net/0b4444914e56a3049a54c16b44f1a6619c0b198e/260e4/images/docs/user-guide/logging/logging-from-application.png)

应用内加入额外的日志处理逻辑，直接发送到 Logging Backend。

## EFK 实践

这一节讲述 EFK 的安装，配置以及使用

#### 安装&配置

可以使用预定义的 Helm Chart 快速安装部署，或者自行单独部署 EFK 各个组件

##### Helm 安装

TODO

##### 自定义安装

TODO 



ref

- [Fluent Bit官方文档](https://docs.fluentbit.io/manual/)
- [K8S日志架构](https://kubernetes.io/docs/concepts/cluster-administration/logging/)
- [K8S 官方 EFK 配置示例](https://github.com/kubernetes/kubernetes/tree/master/cluster/addons/fluentd-elasticsearch)

