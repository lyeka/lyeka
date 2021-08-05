---
title: "Go Channel"
description: 
date: 2021-08-04T14:36:15+08:00
image: 
math: 
license: 
hidden: false
draft: true
categories:
  - program
  - go
tags:
  - go
  - 源码分析
---



## 源码分析

###  内部结构

```go
type hchan struct {
	qcount   uint           // channel中元素个数
	dataqsiz uint           // 环形队列长度
	buf      unsafe.Pointer // 指向dataqsiz的指针
	elemsize uint16 // 元素大小
	closed   uint32 // 是否关闭
	elemtype *_type // 元素类型
	sendx    uint   // 发送索引
	recvx    uint   // 接收索引
	recvq    waitq  // 接收者等待列表
	sendq    waitq  // 发送者等待列表

	lock mutex
}
```

图示

![](https://i.loli.net/2021/08/04/7NGXd9ZI8zrSRxl.png)



### 创建

编译器会将channel的创建转换为`makechan`

```go
// src/runtime/chan.go

// makechan channel的创建
func makechan(t *chantype, size int) *hchan {
	elem := t.elem
	
	// 计算需要分配的内存大小
	// 总共要分配的大小 = mem(元素大小*个数）+ hchanSize
	mem, overflow := math.MulUintptr(elem.size, uintptr(size))
	if overflow || mem > maxAlloc-hchanSize || size < 0 {
		panic(plainError("makechan: size out of range"))
	}

	// 分配内存
	// 包含三种状况
	// 1. mem==0: 元素大小为0或者无缓冲
	// 2. 元素不包含指针
	// 3. 元素包含指针
	var c *hchan
	switch {
	case mem == 0:
		// Queue or element size is zero.
		c = (*hchan)(mallocgc(hchanSize, nil, true))
		// Race detector uses this location for synchronization.
		c.buf = c.raceaddr()
	case elem.ptrdata == 0:
		// Elements do not contain pointers.
		// Allocate hchan and buf in one call.
		c = (*hchan)(mallocgc(hchanSize+mem, nil, true))
		c.buf = add(unsafe.Pointer(c), hchanSize)
	default:
		// Elements contain pointers.
		c = new(hchan)
		c.buf = mallocgc(mem, elem, true)
	}

	c.elemsize = uint16(elem.size)
	c.elemtype = elem
	c.dataqsiz = uint(size)

	return c
}
```

1. 计算所需的内存
2. 创建hchan，分配所需内存
3. 填充其它标识字段

### 发送

TODO



### 接收

TODO



## 使用

TODO





ref

- [Go语言原本 - 3.6 通信原语](https://golang.design/under-the-hood/zh-cn/part1basic/ch03lang/chan/)



