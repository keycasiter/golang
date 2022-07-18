# []()前言

基准测试是测量一个程序在固定工作负载下的性能，Go语言也提供了可以支持基准性能测试的`benchmark`。

# []()使用方法

下面展示一个基准测试的示例代码来剖析下它的使用方式：

```go
func Benchmark_test(b *testing.B) {
	for i := 0; i < b.N ; i++ {
		s := make([]int, 0)
		for i := 0; i < 10000; i++ {
			s = append(s, i)
		}
	}
}
```

* 进行基准测试的文件必须以`*_test.go`的文件为结尾，这个和测试文件的名称后缀是一样的，例如`abc_test.go`
* 参与Benchmark基准性能测试的方法必须以`Benchmark`为前缀，例如`BenchmarkABC()`
* 参与基准测试函数必须接受一个指向`Benchmark`类型的指针作为唯一参数，`*testing.B`
* 基准测试函数不能有返回值
* `b.ResetTimer`是重置计时器，调用时表示重新开始计时，可以忽略测试函数中的一些准备工作
* `b.N`是基准测试框架提供的，表示循环的次数，因为需要反复调用测试的代码，才可以评估性能

# []()命令及参数

性能测试命令为`go test [参数]`，比如`go test -bench=. -benchmem`，具体的命令参数及含义如下：

| 参数            | 含义                                        |
| ------------- | :---------------------------------------- |
| -bench regexp | 性能测试，支持表达式对测试函数进行筛选。                      |
| -bench        | . 则是对所有的benchmark函数测试，指定名称则只执行具体测试方法而不是全部 |
| -benchmem     | 性能测试的时候显示测试函数的内存分配的统计信息                   |
| －count n      | 运行测试和性能多少此，默认一次                           |
| -run regexp   | 只运行特定的测试函数， 比如-run ABC只测试函数名中包含ABC的测试函数   |
| -timeout t    | 测试时间如果超过t, panic,默认10分钟                   |
| -v            | 显示测试的详细信息，也会把Log、Logf方法的日志显示出来            |

# []()测试结果

执行命令后，性能测试的结果展示如下

```bash
$ go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: program/benchmark
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
Benchmark_test-12        7439091               152.0 ns/op           248 B/op          5 allocs/op
PASS
ok      promgram/benchmark    1.304s
```

对以上结果进行逐一分析：

| 结果项               | 含义                                                        |
| :---------------- | :-------------------------------------------------------- |
| Benchmark_test-12 | **Benchmark_test** 是测试的函数名 **-12** 表示GOMAXPROCS（线程数）的值为12 |
| 7439091           | 表示一共执行了**7439091**次，即**b.N**的值                            |
| 152.0 ns/op       | 表示平均每次操作花费了**152.0纳秒**                                    |
| 248B/op           | 表示每次操作申请了**248Byte**的内存申请                                 |
| 5 allocs/op       | 表示每次操作申请了**5**次内存                                         |

# []()性能对比实例

下面通过一个**数字转换字符串**的实例来对比性能测试效果，并进行分析。

```go
//Sprintf
func BenchmarkSprintf(b *testing.B) {
	num := 10
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("%d", num)
	}
}

//Format
func BenchmarkFormat(b *testing.B) {
	num := int64(10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strconv.FormatInt(num, 10)
	}
}

//Itoa
func BenchmarkItoa(b *testing.B) {
	num := 10
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strconv.Itoa(num)
	}
}
```

下面执行命令`go test -bench=. -benchmem`，收到测试报告如下：

```bash
%  go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: program/benchmark
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkSprintf-12     16364854                63.70 ns/op            2 B/op          1 allocs/op
BenchmarkFormat-12      493325650                2.380 ns/op           0 B/op          0 allocs/op
BenchmarkItoa-12        481683436                2.503 ns/op           0 B/op          0 allocs/op
PASS
ok      program/benchmark    4.007s
```

可以发现，**BenchmarkSprintf**方法耗时最长，**BenchmarkFormat**最快，**BenchmarkItoa**也很快。差别在于`fmt.Sprintf()`执行过程中进行了一次内存分配**1 allocs/op**。

# []()结合pprof分析

| 参数                 | 含义          | 命令示例                                             |
| :----------------- | :---------- | :----------------------------------------------- |
| -cpuprofile [file] | 输出cpu性能文件   | `go test -bench=. -benchmem -cpuprofile=cpu.out` |
| -memprofile [file] | 输出mem内存性能文件 | `go test -bench=. -benchmem -memprofile=cpu.out` |

生成的**CPU、内存**文件可以通过`go tool pprof [file]`进行查看，然后在pprof中通过`list [file]`方法查看**CPU、内存**的耗时情况

```bash
### 内存情况
(pprof) list BenchmarkArrayAppend
Total: 36.49GB
ROUTINE ======================== program/benchmark.BenchmarkArrayAppend in /Users/guanjian/workspace/go/program/benchmark/benchmark_test.go
   11.98GB    11.98GB (flat, cum) 32.83% of Total
         .          .      7://Array
         .          .      8:func BenchmarkArrayAppend(b *testing.B) {
         .          .      9:   for i := 0; i < b.N; i++ {
         .          .     10:           var arr []int
         .          .     11:           for i := 0; i < 10000; i++ {
   11.98GB    11.98GB     12:                   arr = append(arr, i)
         .          .     13:           }
         .          .     14:   }
         .          .     15:}
```

```bash
## CPU情况
(pprof) list BenchmarkArrayAppend
Total: 8.86s
ROUTINE ======================== program/benchmark.BenchmarkArrayAppend in /Users/guanjian/workspace/go/program/benchmark/benchmark_test.go
      10ms      640ms (flat, cum)  7.22% of Total
         .          .      6:
         .          .      7://Array
         .          .      8:func BenchmarkArrayAppend(b *testing.B) {
         .          .      9:   for i := 0; i < b.N; i++ {
         .          .     10:           var arr []int
      10ms       10ms     11:           for i := 0; i < 10000; i++ {
         .      630ms     12:                   arr = append(arr, i)
         .          .     13:           }
         .          .     14:   }
         .          .     15:}
```

# []()总结

go提供了benchmark性能测试的工具，提供了对函数使用内存、CPU等情况的报告分析，还可以借助pprof获得更好的分析报告等，如果想要深入分析，还可以使用之前介绍的gdb进行底层代码的链路跟踪，以及对代码进行反编译查看具体的性能损耗情况。

# []()参考

[go benchmark 性能测试 单元测试 基准测试 使用方法详解](https://blog.csdn.net/yzf279533105/article/details/94016601)\
[go benchmark 性能测试](https://www.cnblogs.com/bergus/articles/go-benchmark-xing-neng-ce-shi.html)
