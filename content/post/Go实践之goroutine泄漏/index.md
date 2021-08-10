---
title: "Go实践之goroutine泄漏"
description: 
date: 2021-08-10T19:30:41+08:00
image: 
math: 
license: 
hidden: false
draft: true
category:
    -: go
tag:
    -: TODO
---

## Goroutine泄漏常见场景

- 死循环
- channel阻塞
- 死锁

总结：只要有可能发生死锁的地方的地方就有可能发生goroutine泄漏

### example

TODO



## 如何监控&排查Goroutine泄漏

- `runtime.NumGoroutine()`

- pprof-goroutine

  



ref

- [goroutine泄露：原理、场景、检测和防范](https://segmentfault.com/a/1190000019644257)
- [跟面试官聊 Goroutine 泄露的 6 种方法，真刺激！](https://segmentfault.com/a/1190000040161853)