---
title: "Mongodb索引机制"
description: 
date: 2021-07-09T12:26:24+08:00
image: 
math: 
license: 
hidden: false
draft:
categories:
  - 数据库
tags:
  - MongoDB
  - WiredTiger
  - TODO
---

## WiredTiger存储引擎

在MongoDB 3.2版本之后，默认的存储引擎从MMAPv1改成了WiredTiger。

### 相关命令

**查看存储引擎**

`db.serverStatus`

output

```shell
...
 storageEngine: 
   { name: 'wiredTiger',
     supportsCommittedReads: true,
     oldestRequiredTimestampForCrashRecovery: Timestamp(0, 0),
     supportsPendingDrops: true,
     ...
   }
...
```



**指定存储引擎**

启动时指定

`mongod --storageEngine mmapv1`

或者修改配置文件

```yaml
storage:
   engine: mmapv1
```

### WiredTiger VS MMAPv1

TODO



### WiredTiger 机制

WT支持B树(B-Tree)以及LSM树(Log-Structured Merge Tree) 作为存储结构。

#### B-Tree VS LSM Tree

B树与LSM的实现原理不在这里赘述，会另开文章（TODO）详细介绍。简单而言，B树的读性能更高而LSM的写性能更好，具体benchmark可以看官方的[WiredTiger Btree vs LSM Benchmark ](https://github.com/wiredtiger/wiredtiger/wiki/Btree-vs-LSM)。

![Pretty pictures](https://github.com/wiredtiger/wiredtiger/wiki/attachments/LSM_btree_Throughput.png)

从结果上看，在读吞吐方面，B树高了一个量级不止，在写吞吐方面，LSM大约只提升10%左右，所以综合而言，还是B树更加优秀和普适。

WT默认的存储结构为B树。



#### B-Tree or B+Tree

网络上关于MongoDB的讨论经常会说到MongoDB的索引采用的是B-Tree，而MySQL（InnoDB) 采用的是 B+Tree。这个结论可以说是对的也可以说是错的。对的地方在于B+Tree本就属于B-Tree，是B-Tree的一个优化版本。毕竟MongoDB在v3.2之前确实是B-Tree，哪怕v3.2之后官方文档上也写着

> [1]	MongoDB indexes use a B-tree data structure.

不过在v3.2切换存储引擎为WT后，按WT的设计文档来讲，MongoDB的索引具体来讲实质是B+Tree

> 原文：
>
> WiredTiger maintains a table's data in memory using a data structure called a B-Tree ( B+ Tree to be specific), referring to the nodes of a B-Tree as pages. Internal pages carry only keys. The leaf pages store both keys and values.
>
> 译文：
>
> WiredTiger 使用B-Tree（具体来说是B+Tree）作为数据结构内存中维护数据目录（table），指向以B-Tree作为页结构（page）的节点（node）。
>
> 内部页只存储索引（keys）。叶子页同时存储索引（keys）以及数据 （values）

所以最准确的说法是MongoDB在V3.2之后索引采用B+Tree。

TODO 补充MonogoDB B+Tree索引示意图

## MongDB索引

### 索引类型

支持

- 单字段索引
- 复合（多字段）索引
- 多级类型索引（Multikey Index）
- 地理空间索引（Geospatial Index）
- 文本索引
- 哈希索引

#### 单字段、复合索引

同RMDBS一样，MongoDB支持单字段与多字段索引

```shell
# 单字段索引
db.products.createIndex( { name: 1 } )
# 复合索引
db.products.createIndex(
  { item: 1, quantity: -1 } ,
  { name: "query for inventory" }
)
```

但是MongoDB的索引需要指定顺序（MySQL中不能指定），1为升序-1为降序，MongoDB会按照指定的顺序构建索引树。

显示指定顺序有利于使用复合索引作排序，在某些场景需要其中一些字段降序，另外字段升序排序的话可以使用到。

#### 多级类型索引 

Multikey 索引主要的作用是给数组中的对象建立索引，利用点表示法声明

```shell
db.spu.createIndex( { sku.price: 1 } ))
```

#### 地理空间索引 

支持两种特殊索引

- [2d indexes](https://docs.mongodb.com/manual/core/2d/)
- [2dsphere indexes](https://docs.mongodb.com/manual/core/2dsphere/)

#### 文本索引

支持单字段以及多字段文本索引，但是一个Collection只能有一个文本索引。

```shell
db.reviews1.createIndex( { comments: "text" } )
db.reviews2.createIndex(
   {
     subject: "text",
     comments: "text"
   }
 )
```

支持

1. 包含搜索列表中的任一文本

   ```shell
   db.stores.find( { $text: { $search: "java coffee shop" } } )
   ```

2. 短语匹配

   ```shell
   db.stores.find( { $text: { $search: "java \"coffee shop\"" } } )
   ```

3. 词语排查

     ```shell
     db.stores.find( { $text: { $search: "java shop -coffee" } } )
     ```

4. 匹配度排序

   ```shell
   db.stores.find(
      { $text: { $search: "java coffee shop" } },
      { score: { $meta: "textScore" } }
   ).sort( { score: { $meta: "textScore" } } )
   ```

​    总体看来，和MySQL的全文搜索能力相差无几，需要更加复杂的功能还得上ES等其它外部引擎。

ps: v3.2后支持中文。

#### 哈希索引

为了支持基于哈希的分片，MongoDB提供了哈希索引

```shell
db.collection.createIndex( { _id: "hashed" } )
```

不支持多个字段的复合哈希索引，但在v4.4后支持复合索引中带单个字段哈希索引

```shell
db.collection.createIndex( { "fieldA" : 1, "fieldB" : "hashed", "fieldC" : -1 } )
```

哈希索引只支持等值匹配，无法利用其范围查询。

### 索引属性

支持

- 部分索引
- 稀疏索引
- TTL索引
- 隐藏索引

TODO







ref 

- [WiredTiger Doc](https://source.wiredtiger.com/3.0.0/tune_page_size_and_comp.html)
- [WiredTiger Btree vs LSM Benchmark ](https://github.com/wiredtiger/wiredtiger/wiki/Btree-vs-LSM)
- [MongoDB: How can I change engine type (from B-Tree to LSM-Tree) of _id_ index?](https://stackoverflow.com/questions/59751187/mongodb-how-can-i-change-engine-type-from-b-tree-to-lsm-tree-of-id-index)
- [MongoDB Doc - Index](https://docs.mongodb.com/manual/indexes/)
- [WiredTiger存储引擎之一：基础数据结构分析](https://mongoing.com/topic/archives-35143)

