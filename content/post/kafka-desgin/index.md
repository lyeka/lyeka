---
title: "Kafka Desgin"
description: 
date: 2021-07-21T00:36:13+08:00
image: 
math: 
license: 
hidden: false
draft: true
categories:
  - 消息队列
tags:
  - Kafka
---

## 架构

![](https://i.loli.net/2021/07/21/CRGxzZujSAQDa1H.png)



- producer: 生产消息的业务方
- consumer: 处理消息的业务方
- consumer group: 若干个处理相同类型消息的consumer组成一个consumer group
- topic: 自定义消息的类别，将消息进行逻辑上的划分
- broker: kafka集群实例
- partition: 将topic从物理上划分为若干个区，每个分区在物理上对应一个文件夹

## 持久化

kafka对于ack的消息不作删除处理，而是记录每个消费组ack这个topic的游标



### 文件存储机制

**partition**

对于每个partition，kafka都会创建一个目录来存储其内容，命名规则为topic+parttion序号

**segment**

kfaka会将同个partition的数据拆分成多个大小相等的segment以文件的形式存储，命名格式为offset值，后缀为log，数据将按顺序写入segment文件（kafka提供partition级别的消息顺序保证）。

将数据拆分成多个文件的好处是，但是当消息过期了（kafka支持消息级别的过期机制），且当某个segment最后一条消息过期了(这意味着整个segment内的消息都过期了），可以直接删除这个segment文件，从而避免在一个大文件中为了删除数据而作多余的数据搬迁工作。

**index**

 除了写入数据文件以外，kafka还会生成索引文件，索引文件（index)与数据文件（segment）以对等的形式存在，命名一致，只是后缀名不同，如下图。

![image](https://awps-assets.meituan.net/mit-x/blog-images-bundle-2015/69e4b0a6.png)

每个index文件都存放其对等的log文件的稀疏索引，标记了消息在数据文件的物理偏移地址。如下图：

![image](https://awps-assets.meituan.net/mit-x/blog-images-bundle-2015/c415ed42.png)

N表示segment的第N条消息（注意这并不是消息的offset，而是对应segment的offset），position表示这条消息的物理偏移地址。

ref

- [Redis、Kafka 和 Pulsar 消息队列对比](https://mp.weixin.qq.com/s?__biz=MzAwNTQ4MTQ4NQ==&mid=2453575807&idx=1&sn=9c8875e6800854ed972d323d25c4a5b5&chksm=8cd1cd1dbba6440bf6a31764894840ccd54388e6f74d06e1304a45a36bbc65f638d3899d24a8&exptype=unsubscribed_card_3001_article_onlinev2_3000w_promotion_level1&expsessionid=1884311222645211136&scene=169&subscene=10000&sessionid=1621943446&clicktime=1621943449&enterid=1621943449#rd)
- [Kafka文件存储机制那些事](https://tech.meituan.com/2015/01/13/kafka-fs-design-theory.html)