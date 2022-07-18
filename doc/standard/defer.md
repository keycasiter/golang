# []()数据结构

**defer**的数据结构定义在`$GOROOT/src/runtime/runtime2.go`

```go
// 大体定义如下，忽略少部分字段
type _defer struct {
	sp uintptr //函数栈指针
	pc uintptr //程序计数器
	fn *funcval //函数地址
	link *_defer //指向自身结构的指针， 用于链接多个defer
}
```

# []()规则约定

* 规则一： 延迟函数的参数在defer语句出现时就已经确定
* 规则二： 延迟函数执行按后进先出顺序执行， 即先出现的defer最后执行
* 规则三： 延迟函数可能操作主函数的具名返回值

# []()实现原理

## []()案例

先来看一个简单的demo代码，在`TestDefer()`方法中声明了一个`defer`函数，我们在defer函数前、中、后都对值`s`进行了修改操作，下面通过这个小案例看看它会产生什么东西以及它内部的执行和实现逻辑是怎样的。

```go
5	func main() {
6		s := TestDefer()
7
8		fmt.Println("result => ",s)
9	}
10
11	func TestDefer()(result string)  {
12		s := "init"
13
14		defer func(s string) {
15			s = "defer2"
16		}(s)
17
18		defer func(s string) {
19			s = "defer1"
20		}(s)
21
22		s = "return"
23
24		return s
25	}
```

通过`go tool objdump -s ‘main’ [program]`来查看代码的反编译情况，如下：

```go
TEXT main.TestDefer(SB) /Users/guanjian/workspace/go/program/defer/defer.go
//省略......                       
  defer.go:14           0x10cbd75               e8a6b3f6ff                      CALL runtime.deferprocStack(SB)         
//省略......                         
  defer.go:18           0x10cbdc1               e85ab3f6ff                      CALL runtime.deferprocStack(SB)         
//省略......                             
  defer.go:24           0x10cbe00               e83bbbf6ff                      CALL runtime.deferreturn(SB)            
//省略......                                    
  defer.go:18           0x10cbe16               e825bbf6ff                      CALL runtime.deferreturn(SB)            
//省略......                                   
  defer.go:14           0x10cbe2c               e80fbbf6ff                      CALL runtime.deferreturn(SB)            
//省略......                       
  defer.go:14           0x10cbe40               c3                              RET                                     
  defer.go:11           0x10cbe41               e8fa11faff                      CALL runtime.morestack_noctxt(SB)       
  defer.go:11           0x10cbe46               e995feffff                      JMP main.TestDefer(SB)                  
//省略......      
```

## []()defer初始

`runtime.deferprocStack` 在声明**defer**处调用， 将defer函数存入goroutine的链表中，这相当于defer的**初始阶段**，若defer声明中包含参数，那么参数将被初始化且不会受到后续变更影响，也就是说defer的初始化阶段入参只能初始化赋值。

## []()defer执行

`runtime.deferreturn` 在**return**指令， 准确的讲是在ret指令前调用，将**defer**从goroutine链表中取出并执行，这相当于defer的**执行阶段**。

通过反编译文件来看，是按照如下顺序执行的：
<div align=center>


<img src="https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/fc512dcf643c489a889a644065901f54~tplv-k3u1fbpfcp-zoom-1.image?" alt="314e8ad3b2dcc27a0da50748e2e4f184.jpeg" width="100%" />

</div>

> **多个defer的执行顺序 :**\
> \
> 通过梳理可知defer执行顺序是根据声明顺序的倒序执行的，即声明顺序为**defer-1、defer-> > 2、defer-3**，那么执行顺序则为**defer-3、defer-2、defer-1**，他们是通过`_defer`的`link *_defer`进行串联的。

# []()案例分析

下面有四个方法，**test\_01、test\_02是值传递，test\_03、test\_04是引用传递**\
我们分别从**作用域、栈执行顺序、值传递**等来分析下

## []()案例-1

```go
func main() {
	s1 := test_01()

	fmt.Println("test_01 => ", s1)	//test_01 =>  0
}

func test_01() int {
	var i int
	defer func() {
		fmt.Println("defer-2 => ",i)//这里打印1，上一个defer操作生效了
		i++	//不影响返回值
	}()
	defer func() {
		fmt.Println("defer-1 => ",i)//这里打印0
		i++	//不影响返回值
	}()
	return i //返回值以此为准
}
```

<div align=center>

<img src="https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/8a579e8cce4c4ff0bc706a995665b26e~tplv-k3u1fbpfcp-zoom-1.image" alt="在这里插入图片描述" width="100%" />

</div>

## []()案例-2

```go
func main() {
	s2 := test_02()
	
	fmt.Println("test_02 => ", s2) //test_02 =>  2
}

func test_02() (res int) {
	var i int
	defer func() {
		fmt.Println("defer-2 => ",res)//这里打印1，上一个defer操作生效了
		res++	//影响返回值，最终返回值
	}()
	defer func() {
		fmt.Println("defer-1 => ",res)//这里打印0
		res++	//影响返回值
	}()
	return i	//影响返回值，但并不是唯一和最终影响返回值的地方
}
```

<img src="https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/bc10283db1cd4452b0df41b499d0921a~tplv-k3u1fbpfcp-zoom-1.image" alt="在这里插入图片描述" width="100%" />

## []()案例-3

```go
func main() {
	s3 := test_03()
	fmt.Println("test_03 => ", s3) //test_03 =>  &{jack}
}

type User struct {
	Name string
}

func test_03() *User {
	user := &User{}
	defer func() {
		user.Name="jack"	//不影响返回值
	}()
	return user //返回值以此为准
}
```

<img src="https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/b33306dbf63a4906b82a91b2fe1358ca~tplv-k3u1fbpfcp-zoom-1.image" alt="在这里插入图片描述" width="100%" />

## []()案例-4

```go
func main() {
	s4 := test_04()

	fmt.Println("test_04 => ", s4) //test_04 =>  &{jack}
}

type User struct {
	Name string
}

func test_04() (resUser *User) {
	user := &User{}
	defer func() {
		resUser.Name="jack" //影响返回值，返回值以此为准
	}()
	return user	//影响返回值，但并不是唯一和最终影响返回值的地方
}
```

**案例4**和**案例3**的逻辑基本一样，可参考**案例3**

# []()总结

* defer定义的延迟函数参数在defer语句出时就已经确定下来了
* defer定义顺序与实际执行顺序相反
* return不是原子操作， 执行过程是: 保存返回值(若有)—>执行defer（ 若有） —>执行ret跳转
* 申请资源后立即使用defer关闭资源是好习惯

# []()参考

《Go专家编程》
