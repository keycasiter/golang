---
theme: arknights
---
# 前言
在`Golang`的并发工具中提供了一些非常有意思的实现，本篇来一起探究下非常简单易懂的`sync.Once`。

# 源码分析
如下是`sync.Once`全部源码，结构体声明、方法实现、注释解释加起来一共70多行，非常简短精练，但是却提供了一个在整个对象生命周期内只被调用一次的功能承诺和实现。
```go
// Once对象确保只会被执行一次
// Once对象在首次使用后不要复制.
type Once struct {
   // done 表示是否完成执行
   // 它是结构中的第一个，因为它用于热路径
   // 热路径内联在每个调用点
   // 在某些架构（amd64/386）上，先完成后执行允许更紧凑的指令，而在其他架构上允许更少的指令（用于计算偏移量）。
   done uint32
   m    Mutex
}

// 当且仅当Do在此实例中第一次被调用时，Do调用函数f，换句话说，给定once变量一次值。
// 如果once.Do（f）被多次调用，只有第一次才会真正调用f，即使每次调用中f的值不同。
// 每个函数都需要一个新的Once实例来执行。
// Do用于必须只运行一次的初始化，循环调用once.Do(f)会导致死锁
// 如果f出现panic，once.Do(f)会认为它已经返回了，之后的Do()调用直接返回而不是再调用f函数
func (o *Once) Do(f func()) {
   // 提示: 这是一个错误的Do实现:
   //
   // if atomic.CompareAndSwapUint32(&o.done, 0, 1) {
   //    f()
   // }
   
   // Do 要保证它返回的时候f函数已经完成，而上面这个实现方式不能做到这个保证:
   // 如果进行两次一样的调用, 竞争胜出的来调用f函数, 并且会立即返回不会等待第一次调用完成，
   // 这就是为什么慢路径回落到Mutex，以及为什么atomic.StoreUint32必须延迟到f返回之后

   if atomic.LoadUint32(&o.done) == 0 {
      // 概述了慢速路径（slow-path）以允许快速路径（fast-path）的内联.
      o.doSlow(f)
   }
}

func (o *Once) doSlow(f func()) {
   o.m.Lock()
   defer o.m.Unlock()
   if o.done == 0 {
      defer atomic.StoreUint32(&o.done, 1)
      f()
   }
}
```
## 结构体
### Once

## 方法
### func Do(f func())
### func doSlow(f func())

## 原理提炼

## 坑点总结