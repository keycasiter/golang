# []()数据结构

slice的定义在`$GOROOT/src/runtime/slice.go`

```go
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}
```

`array`指针指向底层数组， `len`表示切片长度， `cap`表示底层数组容量

# []()slice创建

## []()通过make创建

```go
	//make
	slice := make([]int, 5, 10)
```

<img src="https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/e90b29e4968e4c9691fea5a9c8f2d506~tplv-k3u1fbpfcp-zoom-1.image" alt="在这里插入图片描述" width="100%" />

## []()通过数组创建

```go
	//array
	array := [10]int{}
	slice := array[0:5]
```

<img src="https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/946c0a434aef4dc5a674cf2d1935ba32~tplv-k3u1fbpfcp-zoom-1.image" alt="在这里插入图片描述" width="100%" />

## []()内存共享

当**slice**通过**数组**切分时，两者会共用内存空间，此时`slice[0] == array[5] : true
slice[1] == array[6] : true`，这个特性需要特别注意，尤其是在同时处理**数组**和**slice**的过程中，如我们操作`array[5] = 8`，那么`slice[0]`此时也是`8`

<img src="https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/974f16e2ae0548feb5eb0c07e9b1b2d5~tplv-k3u1fbpfcp-zoom-1.image" alt="在这里插入图片描述" width="100%" />

当我们使用`make`方式进行切片初始化的时候经过了哪些处理呢？

```go
	//make
	slice := make([]int, 5, 10)
```

## []()slice初始化

通过`gdb`断点可以看到，使用到了`slice.go`文件中的`makeslice()`方法，如下：

```go
func makeslice(et *_type, len, cap int) unsafe.Pointer {
	mem, overflow := math.MulUintptr(et.size, uintptr(cap))
	if overflow || mem > maxAlloc || len < 0 || len > cap {
		// NOTE: Produce a 'len out of range' error instead of a
		// 'cap out of range' error when someone does make([]T, bignumber).
		// 'cap out of range' is true too, but since the cap is only being
		// supplied implicitly, saying len is clearer.
		// See golang.org/issue/4085.
		mem, overflow := math.MulUintptr(et.size, uintptr(len))
		if overflow || mem > maxAlloc || len < 0 {
			panicmakeslicelen()
		}
		panicmakeslicecap()
	}
	//以上是对内存溢出情况对panic处理
	return mallocgc(mem, et, true)
}
```

# []()slice扩容

slice扩容的方法定义在`$GOROOT/src/runtime/slice.go`的`growslice`方法中。

## []()通用扩容策略

```go
	newcap := old.cap
	doublecap := newcap + newcap
	if cap > doublecap {
		newcap = cap
	} else {
		if old.cap < 1024 {
			newcap = doublecap
		} else {
			// Check 0 < newcap to detect overflow
			// and prevent an infinite loop.
			for 0 < newcap && newcap < cap {
				newcap += newcap / 4
			}
			// Set newcap to the requested cap when
			// the newcap calculation overflowed.
			if newcap <= 0 {
				newcap = cap
			}
		}
	}
```

* 如果新cap大小是当前cap的**2倍**以上，那么按照新cap进行扩容
* cap小于1024，按照**2倍**扩容
* cap大于1024，按照**1.25倍**扩容

通过代码来看下**slice普通扩容过程中len、cap以及内存分配情况**，如下：

```go
// 普通扩容情况，这里是int32类型
func slice() {
	slice := make([]int32, 0)
	for i := 0; i < 10; i++ {
		fmt.Printf("seq=%v, len=%v, cap=%v，\t ptr=%p \t slice=%#v \n",
			i,
			len(slice),
			cap(slice),
			&slice,
			slice)
		slice = append(slice, int32(i))
	}
}
```

输出日志如下：

```bash
seq=0, len=0, cap=0，    ptr=0xc00011a018        slice=[]int32{} 
seq=1, len=1, cap=2，    ptr=0xc00011a018        slice=[]int32{0} 
seq=2, len=2, cap=2，    ptr=0xc00011a018        slice=[]int32{0, 1} 
seq=3, len=3, cap=4，    ptr=0xc00011a018        slice=[]int32{0, 1, 2} 
seq=4, len=4, cap=4，    ptr=0xc00011a018        slice=[]int32{0, 1, 2, 3} 
seq=5, len=5, cap=8，    ptr=0xc00011a018        slice=[]int32{0, 1, 2, 3, 4} 
seq=6, len=6, cap=8，    ptr=0xc00011a018        slice=[]int32{0, 1, 2, 3, 4, 5} 
seq=7, len=7, cap=8，    ptr=0xc00011a018        slice=[]int32{0, 1, 2, 3, 4, 5, 6} 
seq=8, len=8, cap=8，    ptr=0xc00011a018        slice=[]int32{0, 1, 2, 3, 4, 5, 6, 7} 
seq=9, len=9, cap=16，   ptr=0xc00011a018        slice=[]int32{0, 1, 2, 3, 4, 5, 6, 7, 8} 
```

> 日志解释：
>
> * `seq`是执行次序
> * `len`是当前已使用空间
> * `cap`是当前全部容量
> * `ptr`是切片的指针
> * `slice`是切片的内容

借助`benchmark`来查看下内存分配情况：

```bash
 %  go test -bench=SliceExpand -benchmem
goos: darwin
goarch: amd64
pkg: program/slice
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkSliceExpand-12          6996427               144.7 ns/op           248 B/op          5 allocs/op
PASS
ok      program/slice        1.195s
```

**5 allocs/op**表明10次循环过程中进行了5次的内存分配，其实这便是**cap**的扩容过程，即`0 -> 1 -> 2 -> 4 -> 8 -> 16`的5次扩容的内存操作。

## []()特殊扩容策略

对于一些特殊类型，出于内存对齐充分利用的考虑，扩容过程中需要进行特殊处理，下面是特殊处理扩容的策略代码，其中最主要的是`roundupsize()`方法，它在本地存储了各长度的内存对其策略，根据type类型的size来指定扩容情况，这样是对内存友好的。

```go
// Specialize for common values of et.size.
	// For 1 we don't need any division/multiplication.
	// For sys.PtrSize, compiler will optimize division/multiplication into a shift by a constant.
	// For powers of 2, use a variable shift.
	switch {
	case et.size == 1:
		lenmem = uintptr(old.len)
		newlenmem = uintptr(cap)
		capmem = roundupsize(uintptr(newcap))
		overflow = uintptr(newcap) > maxAlloc
		newcap = int(capmem)
	case et.size == sys.PtrSize:
		lenmem = uintptr(old.len) * sys.PtrSize
		newlenmem = uintptr(cap) * sys.PtrSize
		capmem = roundupsize(uintptr(newcap) * sys.PtrSize)
		overflow = uintptr(newcap) > maxAlloc/sys.PtrSize
		newcap = int(capmem / sys.PtrSize)
	case isPowerOfTwo(et.size):
		var shift uintptr
		if sys.PtrSize == 8 {
			// Mask shift for better code generation.
			shift = uintptr(sys.Ctz64(uint64(et.size))) & 63
		} else {
			shift = uintptr(sys.Ctz32(uint32(et.size))) & 31
		}
		lenmem = uintptr(old.len) << shift
		newlenmem = uintptr(cap) << shift
		capmem = roundupsize(uintptr(newcap) << shift)
		overflow = uintptr(newcap) > (maxAlloc >> shift)
		newcap = int(capmem >> shift)
	default:
		lenmem = uintptr(old.len) * et.size
		newlenmem = uintptr(cap) * et.size
		capmem, overflow = math.MulUintptr(et.size, uintptr(newcap))
		capmem = roundupsize(capmem)
		newcap = int(capmem / et.size)
	}

// Returns size of the memory block that mallocgc will allocate if you ask for the size.
func roundupsize(size uintptr) uintptr {
	if size < _MaxSmallSize {
		if size <= smallSizeMax-8 {
			return uintptr(class_to_size[size_to_class8[divRoundUp(size, smallSizeDiv)]])
		} else {
			return uintptr(class_to_size[size_to_class128[divRoundUp(size-smallSizeMax, largeSizeDiv)]])
		}
	}
	if size+_PageSize < size {
		return size
	}
	return alignUp(size, _PageSize)
}
```

通过代码来看下**slice特殊扩容过程中len、cap以及内存分配情况**，如下：

```go
// 特殊扩容情况，这里是int8类型
func slice() {
	slice := make([]int8, 0)
	for i := 0; i < 10; i++ {
		fmt.Printf("seq=%v, len=%v, cap=%v，\t ptr=%p \t slice=%#v \n",
			i,
			len(slice),
			cap(slice),
			&slice,
			slice)
		slice = append(slice, int8(i))
	}
}
```

输出日志如下：

```bash
seq=0, len=0, cap=0，    ptr=0xc0000a8018        slice=[]int8{} 
seq=1, len=1, cap=8，    ptr=0xc0000a8018        slice=[]int8{0} 
seq=2, len=2, cap=8，    ptr=0xc0000a8018        slice=[]int8{0, 1} 
seq=3, len=3, cap=8，    ptr=0xc0000a8018        slice=[]int8{0, 1, 2} 
seq=4, len=4, cap=8，    ptr=0xc0000a8018        slice=[]int8{0, 1, 2, 3} 
seq=5, len=5, cap=8，    ptr=0xc0000a8018        slice=[]int8{0, 1, 2, 3, 4} 
seq=6, len=6, cap=8，    ptr=0xc0000a8018        slice=[]int8{0, 1, 2, 3, 4, 5} 
seq=7, len=7, cap=8，    ptr=0xc0000a8018        slice=[]int8{0, 1, 2, 3, 4, 5, 6} 
seq=8, len=8, cap=8，    ptr=0xc0000a8018        slice=[]int8{0, 1, 2, 3, 4, 5, 6, 7} 
seq=9, len=9, cap=16，   ptr=0xc0000a8018        slice=[]int8{0, 1, 2, 3, 4, 5, 6, 7, 8} 
```

借助`benchmark`来查看下内存分配情况：

```bash
% go test -bench=SliceExpand -benchmem                       
goos: darwin
goarch: amd64
pkg: program/slice
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkSliceExpand-12         25947428                47.71 ns/op           24 B/op          2 allocs/op
PASS
ok      program/slice        2.259s
```

**2 allocs/op**表明10次循环过程中进行了2次的内存分配，其实这便是**cap**的扩容过程，即`0 -> 8 -> 16`的2次扩容的内存操作。

## []()小结

* 切片的**cap**一般处理则按照**2倍**扩容，特殊处理按照`roundupsize`函数扩容，按照特殊处理的cap扩容减少了内存操作次数
* 切片的指针没有发生变化，一直是在同一个数组下进行操作的

# []()slice特殊用法

可以使用如下格式进行切片的使用和截取

| 语法                       | 示例                                                       |
| :----------------------- | :------------------------------------------------------- |
| make[type, len, cap]     | sliceA := make([]int, 5, 10) //length = 5; capacity = 10 |
| slice[start : end]       | sliceB := sliceA[0:5] //length = 5; capacity = 10        |
| slice[start : ]          | sliceC := sliceA[0:] //length = 5; capacity = 10         |
| slice[: end ]            | sliceD := sliceA[:5] //length = 5; capacity = 10         |
| slice[start : end : cap] | sliceE := sliceA[0:5:5] //length = 5; capacity = 5       |

# []()总结

* 创建切片时可跟据实际需要预分配容量， 尽量避免追加过程中扩容操作， 有利于提升性能；
* 切片拷贝时需要判断实际拷贝的元素个数
* 谨慎使用多个切片操作同一个数组， 以防读写冲突

# []()参考

《Go专家编程》\
[Go slice扩容深度分析](https://www.jianshu.com/p/96db3f5c0b0e)
