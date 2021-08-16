---
title: "分布式锁设计与实现"
description: 
date: 2021-08-16T14:14:07+08:00
image: 
math: 
license: 
hidden: false
draft: true
category:
    -: nil
tag:
    -: TODO
---



## 

## 理论设计

在单进程中，我们常常会使用（本地）锁来实现对资源的保护，防止多个线程之间的竞用导致该资源被不正常访问。分布式锁主要解决分布式环境下对同一资源的保护，所以需要引入第三方服务来作锁。锁实质上就是一个变量，先获取到锁的访问方对变量设值，设值之后，其它访问方尝试上锁的时候会知道此时资源被占用了而不去修改资源，从而避免了多方同时修改同一资源而导致出错。因此一般的KV存储均可以用于实现分布式锁。



### 选型比对

- 关系型数据库（MySQL）
- Redis
- Zookeeper
- ETCD

|      | MySQL           | Redis | Zookeeper | ETCD |
| ---- | --------------- | ----- | --------- | ---- |
| CAP  | ALL(Default CA) | AP    | CP        | CP   |
| 性能 | 低              | 高    | 中等      | 中等 |
|      |                 |       |           |      |



## Redis实现

### 上锁

`SET <key> <value> NX PX <expire time>`

- key：指的是资源的唯一标识符
- value: 通常使用随机值，从而避免锁被其他方解锁，这样只有上锁的一方才有资格释放锁
- expire time：锁的过期时间，为了避免上锁后系统A进入长时间阻塞或者崩溃，导致该资源长时间或者永久不可被访问



## ETCD实现

TODO



ref

- [Why isn't RDBMS Partition Tolerant in CAP Theorem and why is it Available?](https://stackoverflow.com/questions/36404765/why-isnt-rdbms-partition-tolerant-in-cap-theorem-and-why-is-it-available)
- [基于Redis的分布式锁到底安全吗（上）？](http://zhangtielei.com/posts/blog-redlock-reasoning.html)



