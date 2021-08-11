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

**只要有可能发生阻塞的地方（导致无法正常退出函数）就有可能发生goroutine泄漏**

发生goroutine泄漏可能是程序员本身的业务代码逻辑有误导致的，如上述的场景；

但也有可能是不恰当的使用了库导致，如使用Go的`http.Client`向外部接口发送请求，但没有设置超时时间，假如外部接口崩了迟迟没有响应，那这个请求的goroutine就发生了泄漏。

下面列举一些导致goroutine泄漏的示例，demo可能只是个玩具很可能不会在真实环境这么写，但没关系，重点在于**理解为什么会发生goroutine泄漏**以及**如何监控和排查goroutine泄漏**

### 死循环

```go
func LoopCaseLeak() {
	go func() {
		l := make([]int, 10)
		for len(l) != 0 {
			// ...
		}
	}()
}
```

### Channel阻塞



#### unbuffer channel

**只写不读**

```go
func GoLeak1() {
	ch := make(chan struct{})
	go func() {
		ch <- struct{}{}
	}()
}
```

**只读不写**

```go
func GoLeak2() {
	ch := make(chan struct{})
	go func() {
		<- ch 
	}()
}
```



#### buffered channel

原理同unbuffered channel阻塞导致goroutine泄漏是一样的，不过因为有buff解耦了读写的阻塞关系，所以阻塞发生相对隐蔽

**写满了，但是没有读**

```
func bufferChannelFilled(ch chan int) {
   for i:=0; i<100; i++ {
      ch <- i
   }
}
```

**通过range接收channel，但channel没关闭**

```go
func rangeCauseLeak(ch chan int) {
	for v := range ch {
		fmt.Println(v)
	}
}
```



#### nil channel

对于一个nil channel（通常出现在已经声明但没有初始化的channel变量），对于其的读写都会发生阻塞行为

```go
// 假设下面function传入的都是nil channel

func nilChannelRead(ch chan struct{}) {
	<- ch
}

func nilChannelWrite(ch chan struct{}) {
	ch <- struct{}{}
}

```



### 外部接口阻塞

```go
func CallHttpClient()  {
	_, err := http.Get("https://www.v2ex.com/")
	if err != nil {
		fmt.Println("get error: ", err)
		return
	}
}
```



### 死锁

**多次上锁但是没有释放**

```go
func lockToLeak() {
	var m sync.Mutex
	for i:=0; i<5; i++ {
		go func() {
			m.Lock()
			// ...
		}()
	}
}

```

**WaitGroup 忘记Done或者Add与Done的次数不相等**

```go
func waitGroupToLeak() {
	var wg sync.WaitGroup
	for i:=0; i<10; i++ {
		go func() {
			wg.Add(1)
			// ... do something
			// forget wg.Done
			// wg.Done()
		}()
	}
	time.Sleep(1*time.Second)
	wg.Wait()
}
```



## 如何监控&排查Goroutine泄漏

业务代码出差难以避免，更何况Go中创建一个goroutine太方便了，所以如何监控和排查goroutine泄漏很关键。

导致泄漏的示例代码，主要是上面列举的goroutine泄漏示例集合。

```go
package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"time"
)

func main() {
	go http.ListenAndServe("0.0.0.0:6060", nil)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()
	
	go func() {
		for {
			time.Sleep(2 * time.Second)
			fmt.Println("goroutine nums: ", runtime.NumGoroutine())
		}
	}()


	go func() {
		for {
			time.Sleep(2 * time.Second)
			go LoopCaseLeak()
		}
	}()

	go func() {
		for {
			time.Sleep(2 * time.Second)
			go GoLeak1()
		}
	}()

	go func() {
		for {
			time.Sleep(2 * time.Second)
			go GoLeak2()
		}
	}()

	go func() {
		for {
			time.Sleep(2 * time.Second)
			ch := make(chan int, 10)
			go bufferChannelFilled(ch)
		}
	}()

	go func() {
		for {
			time.Sleep(2 * time.Second)
			ch := make(chan int, 10)
			go rangeCauseLeak(ch)
		}
	}()

	go func() {
		for {
			time.Sleep(2 * time.Second)
			var nilChan chan struct{}
			go nilChannelWrite(nilChan)
			go nilChannelRead(nilChan)
		}
	}()

	go func() {
		for {
			time.Sleep(2 * time.Second)
			go CallHttpClient()
		}
	}()

	go func() {
		for {
			time.Sleep(2 * time.Second)
			go lockToLeak()
		}
	}()

	go func() {
		for {
			time.Sleep(2 * time.Second)
			go waitGroupToLeak()
		}
	}()


	ch := make(chan int64)
	for {
		time.Sleep(1 * time.Second)
		n := time.Now().Unix()
		go func() {
			ch <- n
		}()
	}
}

func GoLeak1() {
	ch := make(chan struct{})
	go func() {
		ch <- struct{}{}
	}()
}

func GoLeak2() {
	ch := make(chan struct{})
	go func() {
		<-ch
	}()
}

func LoopCaseLeak() {
	go func() {
		l := make([]int, 10)
		for len(l) != 0 {
			// ...
		}
	}()
}

func CallHttpClient() {
	_, err := http.Get("https://www.v2ex.com/")
	if err != nil {
		fmt.Println("get error: ", err)
		return
	}
}

func rangeCauseLeak(ch chan int) {
	for v := range ch {
		fmt.Println(v)
	}
}

func bufferChannelFilled(ch chan int) {
	for i:=0; i<100; i++ {
		ch <- i
	}
}

func nilChannelRead(ch chan struct{}) {
	<- ch
}

func nilChannelWrite(ch chan struct{}) {
	ch <- struct{}{}
}

func lockToLeak() {
	var m sync.Mutex
	for i:=0; i<5; i++ {
		go func() {
			m.Lock()
			// ...
		}()
	}
}

func waitGroupToLeak() {
	var wg sync.WaitGroup
	for i:=0; i<10; i++ {
		go func() {
			wg.Add(1)
			// ... do something
			// forget wg.Done
			// wg.Done()
		}()
	}
	time.Sleep(1*time.Second)
	wg.Wait()
}


```



### 利用runtime包

可以利用`runtime.NumGoroutine()`输出目前的goroutine数量，如果观察到进程内goroutine数量持续上涨就有可能发生了goroutine泄漏

```go
go func() {
	for {
		time.Sleep(2 * time.Second)
		fmt.Println("goroutine nums: ", runtime.NumGoroutine())
	}
}()
```



### 利用prometheus监控

prometheus提供了[sdk](https://github.com/prometheus/client_golang)用于监控程序的metrics，用于上报各种Go程序指标，可以搭配grafana可视化来监控

```go
go func() {
    http.Handle("/metrics", promhttp.Handler())
    http.ListenAndServe(":2112", nil)
}()
```

这里简单通过`http://localhost:2112/metrics`地址来观测一下Go进程的metrics

下图中的 `go_goroutines`指示了目前的goroutine数量

![image-20210811145850434](https://i.loli.net/2021/08/11/IBk4fhHLjirKugO.png)



搭配grafana后效果图, [图源](https://ms2008.github.io/2019/06/02/golang-goroutine-leak/)

![img](https://ms2008.github.io/img/in-post/goroutine-leak/goroutine-leak-1.png)

可视化后很方便的观察goroutine是否在持续上涨

### pprof



TODO





ref

- [goroutine泄露：原理、场景、检测和防范](https://segmentfault.com/a/1190000019644257)
- [跟面试官聊 Goroutine 泄露的 6 种方法，真刺激！](https://segmentfault.com/a/1190000040161853)
- [INSTRUMENTING A GO APPLICATION FOR PROMETHEUS](https://prometheus.io/docs/guides/go-application/)
- [Goroutine 泄露排查](https://ms2008.github.io/2019/06/02/golang-goroutine-leak/)
- [golang pprof 实战](https://blog.wolfogre.com/posts/go-ppof-practice/)

