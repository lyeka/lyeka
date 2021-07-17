---
title: "Mysql Overview"
description: 
date: 2021-07-14T11:55:13+08:00
image: 
math: 
license: 
hidden: false
draft: true
categories:
  - 数据库
  - Overview
tags:
  - MySQL
---


## 隔离级别有哪些，分别会出现什么情况
| 隔离级别                    | 脏读 | 不可重复读 | 幻读 |
| --------------------------- | ---- | ---------- | ---- |
| 读未提交（read-uncommitted) | 是   | 是         | 是   |
| 读已提交（read-committed)   | 否   | 是         | 是   |
| 可重复读（repeatable-read)  | 否   | 否         | 是   |
| 串行化（serializable)       | 否   | 否         | 否   |


- 脏读：事务没有提交的修改对其他事务也是可见的
- 不可重复读：事务内同样的查询可能会有不同的结果，因为别的事务在此之间做了修改
- 幻读：某个事务在读取某个范围内的记录后，别的事务又在该范围内插入了新的行



ps: InnDB通过MVCC快照读以及next-key lock解决幻读问题，并且在可重复读级别下默认开启next-key lock

---

## MVCC是什么，是如何实现的，算法说一下？

MVCC多版本并发控制，用于在并发事务中减少锁的使用开销，以提供并发性能


### 实现原理

每行数据聚簇索引会包含两个列

- trx_id——一个事务每次对某条聚集索引记录进行改动时，都会把该事务的事务id赋值给trx_id隐藏列
- roll_point——每次对某条聚簇索引记录进行改动时，都会把旧的版本写入undo日志中。这个隐藏列就相当于一个指针，通过他找到该记录修改前的信息



事务开始时会创建一个readview
包括以下几个内容

- m_ids：在生成ReadView时，当前系统中活跃的事务id列表
- min_trx_id：在生成ReadView时，当前系统中活跃的最小的事务id，也就是m_ids中的最小值
- max_trx_id：在生成ReadView时，系统应该分配给下一个事务的事务id值
- creator_trx_id：生成该ReadView的事务的事务id

ps: max_trx_id并不是m_ids中的最大值，而是全局递增的事务id。比如现在有事务id为1，2，3这三个事务，之后事务id为3的事务提交了，当有一个新的事务生成ReadView时，m_ids的值就包括1和2，min_trx_id的值就是1，max_trx_id的值就是4



然后基于trx_id判断该行的哪个版本对当前事务是可见的
流程图如下
![image.png](https://i.loli.net/2021/07/14/gZJsSUnyc9t5BLb.png)

总结
**trx_id小于max_trx_id且不在m_ids列表的对当前事务可见，不满足这个条件的通过回滚指针将记录回溯到该范围内**


### 如何解决幻读

#### 为什么会出现幻读？

如何A事务只是纯读的话，不会出现幻读的情况，因为selcet使用的的是快照读（当执行select操作是innodb默认会执行快照读，会记录下这次select后的结果，之后select 的时候就会返回这次快照的数据，即使其他事务提交了不会影响当前select的数据，这就实现了可重复读了）；
而当A事务内发生可更新行为的话，此时的读是当前读，在执行这几个操作时会读取最新的版本号记录，写操作后把版本号改为了当前事务的版本号，所以即使是别的事务提交的数据也可以查询到；


#### 如何解决幻读？

1. 在相应的列上建非唯一索引
2. `select ... from ... for update`

next-key lock 相当于 index-record lock 加上 gap lock，通过锁住相关范围内的行禁止别的事务在锁范围内插入行来解决幻读问题。

TODO next-key lock 锁定范围解析



#### 结论

**InnoDB使用mvcc解决了快照读下的幻读**
**使用net-key lock解决了当前读下的幻读**



ref

- [MySQL Doc - 14.7.4 Phantom Rows](https://dev.mysql.com/doc/refman/5.7/en/innodb-next-key-locking.html)
- [MySQL Doc - Next-Key Locks](https://dev.mysql.com/doc/refman/5.7/en/innodb-locking.html#innodb-next-key-locks)

------



## InnoDB与Myisam的区

- InnoDB支持事务，外键
- InnoDB的锁粒度支持到行级，Myisam只支持表锁
- InnoDB采用[聚簇索引](https://www.yuque.com/lumingheng/vfbuv2/gi0e3e)
- Myisam保存总行数

---



## Mysql索引用的是什么算法，为什么
使用了B+树作为底层数据结构，
优点

- 效率高，O(logn)效率的范围查找（因为叶子节点使用链表串联起来）且稳定，相对B树来说
- 相对于平衡二叉树来说，树的高度低，减少IO次数（毕竟要存储在磁盘上）
- 除叶子节点外只保存索引，所以节点较小，这样相同的页可存储更多的节点，也就是说每个节点的m值更多，树更低



ref

- [B+树：MySQL数据库索引是如何实现的？](https://time.geekbang.org/column/article/77830)
- [MySQL的索引，为什么是B+而不是平衡二叉树](https://segmentfault.com/a/1190000023402876)

---





## 为什么使用B+数做索引
对比B数

1. B+非叶子节点只存储索引值，节点更小，一页可以存储更多的节点，减少了IO次数
1. B+树的叶子节点使用链表链接起来，更加适合范围查询

对比红黑树（AVL树）

1. B+树是多路树，红黑树是二叉树，相同数据量下，二叉树的高度会更高，IO次数就增加。AVL一般在内存中才作为索引，B+树适合存储在磁盘上

对比哈希（mysql其实也支持哈希索引）

1. 哈希索引不支持模糊查询
1. 哈希索引不支持范围查询

---



## 为什么 MongoDB 索引选择B-树，而 Mysql 索引选择B+树
MongoDB作为NoSQL更普通用于KV查找，数据聚合性更高，B树节点中存储了索引以及数据，单点查找效率更高；
Mysql作为关系型数据库，范围查询需求高，B+树中的叶子节点串联起来非常适合范围查找，且B+数索引节点小的优点也使其更加格式外部存储


ref

- [为什么 MongoDB 索引选择B-树，而 Mysql 索引选择B+树](https://www.cnblogs.com/kaleidoscope/p/9481991.html)

---



## B+树索引的最底层单元是什么？什么决定了B+树的高度？B+树的叶子节点是单向链表还是双向链表?
链表


### 什么决定了B+树的高度

- Innodb_page_size——Innodb每页的大小
- 索引的大小+指针的大小——决定了B+树非叶子节点存放关键字的数量



### B+树的叶子节点是单向链表还是双向链表?
双向，可以用于order desc


ref

- [为什么生产环境中B+树的高度总是3-4层？](https://zhuanlan.zhihu.com/p/86137284)
- [mysql的B+树如何存储主键和数据，磁盘io和innodb页大小的一些问题](https://blog.csdn.net/LJFPHP/article/details/105318995)

---

## 怎么加索引，索引什么情况下失效，联合索引什么时候失效，覆盖索引了解吗？
### 索引什么情况下失效

- 不满足最左前缀法则
    - 使用范围查询后右边的索引列失效
- 在索引列上使用函数
- like以通配符号开头
- 优化器判定使用索引成本比全表扫描成本大



### 联合索引什么时候失效
不满足最左前缀法则

### 覆盖索引
二级索引字段无法满足select字段需要，需要取回主键id到主键索引（聚簇索引）获取所需字段的过程称回表
覆盖索引是指在二级索引上的索引列满足了查询所需要的字段，无需再回表查询


ref

- [MySQL中IS NULL、IS NOT NULL、!=不能用索引？胡扯！](https://juejin.cn/post/6844903921450745863)
- [MySQL 覆盖索引详解](https://juejin.cn/post/6844903967365791752)

---

## 为什么不用uuid做主键，影响的写入性能还是读取性能？如果业务上能保证唯一性，那么还需要建唯一索引吗？会影响写入性能吗？

- 影响写性能
    - 使用UUID作为主键索引插入的话，UUID是无序的，插入数据页也就不是顺序的，需要做额外的工作寻找合适的位置。写入乱序会导致频繁的页分裂操作，也就导致大量数据行的移动。而且频繁的页分裂会使页变得稀疏并被不规则的填充，导致数据碎片
- 读性能也会稍微下降
    - 毕竟UUID一般比较大，而且占用空间也大



ref

- [聚簇索引和二级索引](https://www.yuque.com/gynb7h/bkgq7y/qk8hph)

---

## ACID了解吗，mysql是用什么机制保证的，redo log和undo log说一下
### mysql使用什么机制保证ACID
mysql使用锁实现了隔离性；使用redo log实现了原子性和持久性；使用undo log实现了一致性


### redo log & undo log
#### redo log
**为什么需要redo log?**
对于数据页的修改不是落盘的，内存中和物理中的数据页会存在差异
因为实时落盘性能差

- 小修改需要刷新整个页
- 页不是连续的，实时落盘会产生大量随机IO

因此在落盘间隙间DB挂了就会丢失数据
因为需要使用redo log 来提高性能的同时保证持久性
为什么redo log性能高

- 日志小（物理日志）
- 随机（IO）写入



**redo log工作流程**

1. commit前将所有修改写到redo log buffer
1. fsync将redo log buffer刷新到磁盘（可配置）
1. commit

**
**redo log恢复**
InnoDB在启动之前都会尝试执行恢复操作（不管上次数据库是否正常关闭）
首先比对刷新到磁盘的LSN(Log Sequence Number)和redo log LSN，如果redo log LSN不等于（大于）磁盘上的LSN，则使用redo log 中两LSN范围内的日志进行数据恢复
至于redo log的格式以及是怎么恢复数据的，太复杂，#TODO




#### undo log
保存了事务发生之前的数据的一个版本，可以用于回滚，同时可以提供多版本并发控制下的读（MVCC），也即非锁定读
逻辑格式的日志，在执行undo的时候，仅仅是将数据从逻辑上恢复至事务之前的状态，而不是从物理页面上操作实现的，这一点是不同于redo log的。undo log将生成与更新语句相反的语句来执行回滚（如insert生成delete）


ref

- 《Mysql技术内幕——InnoDB存储引擎》第7章
- [MySQL之Redo Log](https://zhuanlan.zhihu.com/p/35355751)
- [MySQL中的重做日志（redo log），回滚日志（undo log），以及二进制日志（binlog）的简单总结](https://www.cnblogs.com/wy123/p/8365234.html)



---

## 锁分哪几种，行锁是锁在哪里的？


### 锁的分类，基于悲观锁实现

- 表锁 —— 没有命中索引则加表锁
- 行锁（**基于索引**）—— 当命中索引的话加行锁
    - 记录锁（record lock）—— 当命中了唯一索引（包括主键索引）时
    - 间隙锁（gap lock）—— 当命中了非唯一索引，且范围之中没有符合的数据行
        - 只出现在RR隔离级别
        - 锁定查询条件的左开右闭区间范围内的行
    - 临键锁（next-key lock）—— 当命中了非唯一索引，且范围内有符合的数据行
        - 只出现在RR隔离级别
        - 锁定范围为索引B+树上的结点，如果没有命中具体结点，则为无穷



ref

- [深入理解数据库行锁与表锁](https://zhuanlan.zhihu.com/p/52678870)

---

## mysql 主从同步的过程是什么样的


主从同步好处

- 读写分离，提高性能
- 容灾备份，高可用



形式

- 一主一从
- 一主多从
- 级联赋值
- 多主一从
- 双主复制





复制模式

- 异步模式（默认）—— 主节点无需等待从节点工作状态
- 半同步模式——部分从节点commit主节点才commit
- 同步模式——全部从节点commit后主节点才commit



binlog记录格式

- row——记录数据的变化
- statement——记录SQL
- mixed——混合使用，由mysql判断使用的格式



### 同步过程
三个线程

- 主binlog dump 线程
    - 对于每一个从服务器都会创建一个单独的binlog dump线程
- 从IO线程
- 从SQL线程



1. 记录binlog后，主库发送信号量通知唤醒dump线程（假如该线程追赶上了主库，进入睡眠状态）；
1. IO线程通过将dump线程发送的事件记录到中继日志中（relay log）（通过log_pos来确定同步事件范围）
1. SQL线程读取中继日志来执行SQL



ref

- 《高性能MySQL》第十章+复制
- [深度探索MySQL主从复制原理](https://zhuanlan.zhihu.com/p/50597960)

---

## mysql的高可用是怎么做的

- 避免单点故障
- 故障转移和故障恢复

TODO




- 《高性能MySQL》第十二章——高可用

---

## 慢查询怎么处理？explain的type有几种类型，你常见的有哪几种？
### 慢查询优化思路

1. 开启慢日志
1. 针对慢SQL explain 进行优化
1. 对整体查询进行优化，确定是否是因为高并发引起锁的超时问题
    1. 分库分表
    1. 业务优化
4. MySQL服务器优化
    1. 系统参数



### explain type类型

- system
    - 系统表
- const
    - **命中了主键、唯一索引且等值查询**
- eq_ref
    - **连表**查询中**一对一**关系
- ref
    - **连表**查询中**不是一对一**关系
    - 非唯一索引的**等值查询**
- range
    - 索引上的**范围查询**（between, in, <=>）
- index
    - 需要扫描索引列上的全部数据
- all
    - 全表扫描
- fulltext
    - 全文搜索
- index_merge
    - 索引合并
- ref_or_null
    - 与ref相似，但是包含NULL
- unique_subquery
    - value in (select...)， 子查询是唯一的列？
- index_subquery
    - value in (select...)， 子查询不是唯一的列？



ref

- [Mysql Explain之type详解](https://juejin.cn/post/6844904149864169486)
- [MySQL explain type详解](https://blog.51cto.com/lijianjun/1881208)

---

## 你们的mysql数据量有多大，如何分库分表，分库分表以后如何查询，如何做分布式事务
### 如何分库分表

- 垂直（优先）
    - 库
        - 业务表拆分不同的库
    - 表
        - 将表的字段拆分成多张表（相关度高的放在一起）
- 水平
    - 水平分库分表
        - 作分片
            - RANGE（从0到10000一个表，10001到20000一个表）
            - 哈希
            - 业务模块(如地理位置，时间...)



### 分库分表以后如何查询，如何做分布式事务
TODO


ref

- [MySQL 分库分表方案，总结的非常好！](https://juejin.cn/post/6844903648670007310)

---

## 怎么优化mysql（服务器配置）
TODO


ref

- [优化MySQL服务器设置](https://binism.github.io/blog/2019/12/15/MySQL-params/#innodb%E7%BC%93%E5%86%B2%E6%B1%A0buffer-pool)

---

## 了解TiDB吗，TiDB是如何实现的
分布式关系型数据库

- 水平弹性扩展
- 故障自恢复及异地多活
- 一致性的分布式事务
- 高度兼容 MySQL



TODO
