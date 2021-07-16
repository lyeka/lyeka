---
title: "Redis Interview"
description: 
date: 2021-07-15T11:53:51+08:00
image: 
math: 
license: 
hidden: false
draft: true
categories:
    - interview
tag:
    - redis
---

## redis 数据类型

### 外部数据类型

基础类型

- strings
- hashes
- lists
- sets
- sorted sets

高级类型

- bitmaps
- hyperloglogs
- geospatial indexs
- streams

### 内部数据结构

- dict 
- SDS

- ziplist
- quicklist
- skiplist



ref

- [An introduction to Redis data types and abstractions](https://redis.io/topics/data-types-intro)

---

## redis string类型的底层数据结构
```c
struct sdshdr {
    int len;  // buf已用长度
    int free; // buf剩余长度
    char buf[];  // 数据实际存储的地方
}
```

why:

1. 高效计算长度以及修改（追加等）操作
   1. 常数时间复杂度度计算字符串长度
   2. 减少内存分配次数
2. 支持二进制安全
   1. 在 Redis 中， 客户端传入服务器的协议内容、 aof 缓存、 返回给客户端的回复， 等等， 这些重要的内容都是由 sds 类型来保存的。

ref

- [Simple Dynamic String](https://redisbook.readthedocs.io/en/latest/internal-datastruct/sds.html)
- [二进制安全(binary safe)是什么意思？](https://www.zhihu.com/question/28705562)

---

## sorted set实现

基于跳表（skiplist)——多层链表，在底层全量数据的基础上，加上多层稀疏链表构成。查找时自顶向下遍历链表，类似于二分查找，将查找复杂读降低至O(logN)。



### skiplist

插入

![skiplist_insertions.png](https://i.loli.net/2021/07/16/qsGalfT3iXKLxYF.png)

redis实现的skiplist相邻两层并不是2:1的比例，而是通过随机数生成的算法构建多层链表，从而避免插入为了保证2:1比例而作额外的工作。



查找

![search_path_on_skiplist.png](https://i.loli.net/2021/07/16/9ml1Pk3CfJja85E.png)

ref

- [Redis内部数据结构详解(6)——skiplist](http://zhangtielei.com/posts/blog-redis-skiplist.html)

---

## 主从同步原理



为什么作主从

- 备份
- 容灾
- 读写分离，负载均衡

### 同步方式

redis2.8之前，同步方式只有全量复制一种模式，从节点通过`sync`命令向主节点请求全量同步；redis2.8之后支持了`psync`用于部分复制。

主从节点都会保存一份复制偏移量（offset）用于部分复制；

`psync`不仅只支持部分复制，主从节点都可以将某次同步改为全量复制，示意图如下

![img](https://images2018.cnblogs.com/blog/1174710/201806/1174710-20180628011547892-692403928.png)

全量复制的效率比较低，流程如下：

1. 主节点调用`bgsave`fork子进程进行RDB持久化
2. 将RDB文件通过网络发送至从节点
3. 从节点清空旧数据，将接收到的RDB文件载入到内存中
4. 主节点将前面执行前面步骤期间新产生的的复制缓冲区里面数据同步给从节点以追至最新状态



部分复制实现：

1. 主节点在执行写命令的同时还会将其存储一份到复制缓存区（FIFO队列）
2. 主从节点分别维护一份复制偏移量
3. 通过复制偏移量差来取复制缓冲区的数据来部分复制



ref

- [深入学习Redis（3）：主从复制](https://www.cnblogs.com/kismetv/p/9236731.html)

---

## 哨兵原理



