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

- qcount是channel中已经存在的数据的个数，dataqsiz是channel的cap，当qcount=dataqsiz时，channel阻塞，所以qcount会永远小于等于dataqsiz
- channel内部通过互斥锁（mutex）来保证线程安全，社区中有一些无锁channel的实现，当因为无法保证FIFO或者多核条件下的性能等问题没有被接纳



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



```go
// src/runtime/chan.go

func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
    // 空channel的发送会阻塞
	if c == nil {
		if !block {
			return false
		}
		gopark(nil, nil, waitReasonChanSendNilChan, traceEvGoStop, 2)
		throw("unreachable")
	}

	...

	lock(&c.lock)

	// 向关闭的channel发送数据直接panic
	if c.closed != 0 {
		unlock(&c.lock)
		panic(plainError("send on closed channel"))
	}

	// 1. 如果有等待的goroutine队列，从等待的队列中出列(FIFO), 向其发送数据
	if sg := c.recvq.dequeue(); sg != nil {
		send(c, sg, ep, func() { unlock(&c.lock) }, 3)
		return true
	}

	// 2. buff有空间，将数据放入buff
	if c.qcount < c.dataqsiz {
		// Space is available in the channel buffer. Enqueue the element to send.
		// 入buff
		qp := chanbuf(c, c.sendx)
        // 拷贝数据到buff
		typedmemmove(c.elemtype, qp, ep)
		// 更新发送索引，注意这里的buff是个环形队列
		c.sendx++
		if c.sendx == c.dataqsiz {
			c.sendx = 0
		}
		// 元素个数+1
		c.qcount++
		// 解锁返回
		unlock(&c.lock)
		return true
	}

	if !block {
		unlock(&c.lock)
		return false
	}

	// 3. buff无空间，阻塞，让出P使用权
	gp := getg()
	// 创建一个sudog,加入channel的sendq
	mysg := acquireSudog()
	...
	c.sendq.enqueue(mysg)
	// 将发送操作的goroutine陷入阻塞睡眠，直至被chanparkcommit被唤醒
	gopark(chanparkcommit, unsafe.Pointer(&c.lock), waitReasonChanSend, traceEvGoBlockSend, 2)
	
    ...
    
	return true
}
```

发送可以分成三种情况

1. 有等待接收的goroutine：直接发送
2. 没有等待接收的goroutine，缓冲区未满：将数据拷贝至缓存区，发送行为返回
3. 没有等待接收的goroutine，缓存区已满：阻塞，该goroutine让出P



### 接收

TODO



## 使用

TODO





ref

- [Go语言原本 - 3.6 通信原语](https://golang.design/under-the-hood/zh-cn/part1basic/ch03lang/chan/)



