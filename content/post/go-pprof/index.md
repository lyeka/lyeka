---
title: "Go Pprof"
description: 
date: 2021-08-03T16:10:41+08:00
image: 
math: 
license: 
hidden: false
draft: true
category:
    - nil
tag:
    - TODO
---



## 具体profilling

### CPU

命令

```shell
# 利用curl下载profile文件
curl -o cpu-profile http://127.0.0.1:6060/debug/pprof/profile

# 通过终端(命令)分析profile文件
go tool pprof -http=:8888 .\cpu-profile

# 通过http可视化分析
go tool pprof -http=:8888 .\cpu-profile
```

