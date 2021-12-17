---
title: "MongoDB安全机制"
description: 
date: 2021-07-15T14:32:23+08:00
image: 
math: 
license: 
hidden: false
draft: true
category:
    -: 数据库
tag:
    -: MongoDB
---

## 开启安全机制

### 创建超级用户

```shell
use admin
db.createUser(
  {
    user: "root",
    pwd: "<your password>", 
    roles: [ { role: "userAdminAnyDatabase", db: "admin" }, "readWriteAnyDatabase" ]
  }
)
```



### 重启mongo以启用安全机制

**关机**

```shell
db.adminCommand( { shutdown: 1 } )
```

**开启访问控制（开关）**

有两种方式

1. 启动时带上参数

   1. `mongod --auth ...`

2. 修改配置文件开启

   1. ```yaml
      security:
          authorization: enabled
      ```

   2. `mongod`

只有开启了auth，基于角色控制访问的限制才会生效

