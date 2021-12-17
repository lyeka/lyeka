---
title: "接口优化思路"
description: 
date: 2021-12-17T19:50:26+08:00
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

# 接口优化记录

## jemeter压测

### 利用GUI生成测试计划

### 使用命令行压测





## pprof查看性能瓶颈

下载profile结果文件

`curl <host>/debug/pprof/profile -o <profile file name>`

go tool 分析

`go tool pprof <profile file name>`

go tool 可视化

`go tool pprof -http=:8888 <profile file name>`

### CPU