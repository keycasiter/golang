>🔰 全文字数 : 8K+\
>🕒 阅读时长 : 10min\
>📋 关键词汇 : golang / reflect\
>👉 欢迎关注 : [大摩羯先生](https://mp.weixin.qq.com/mp/profile_ext?action=home&__biz=MzkxMDIxMzk5Mw==#wechat_redirect)

# 什么是反射
&emsp;&emsp;这篇主要聊聊`Golang`中的`Reflect`，也就是**反射**。`Golang`是一种强类型、静态类型的语言，在编译期就已经确定好每个变量的类型，**反射提供的是程序在运行时可以访问、检测、修改自身状态或行为的一种能力，使得编程语言能够有一定的动态能力**。

&emsp;&emsp;众所周知，编程语言都是依靠一定的组织结构来构建的，比如**代码块、类、方法、字段**等，而这里面最原子的表达单位就是承载一个个数据传递的变量，变量又会按照特定的数据格式表示为**字符串**、**数值型**、**布尔型**等各种类型，我们可以用字符串描述一段话，用数值来实现数学运算，用布尔表示是或否，在利用编程语言提供的不同类型构建编码时，其实都是对现实世界的一种抽象和映射，而且这一切都是建立在明确的意图表达上的，这种明确的意图表达能否让我们选择编程语言提供的具体的某一个数据类型来匹配。

&emsp;&emsp;特别的是，现实世界里的诉求是复杂的，有很多时候我们并不能提前明确意图，比如我们要实现一个收集用户信息的方法，最开始我们只定义允许用户传入一个字符串来传递姓名即可，如下：

```go
func CollectUserInfoV1(name string) {
   fmt.Println("姓名:", name)
}
```
随着需求的丰富，我们还要收集用户的年龄、住址、电话等信息，于是我们可以把这些信息封装到一个结构体中，直接传递这个结构体，如下：
```go
func CollectUserInfoV2(user *User) {
   fmt.Println("姓名:", user.Name)
   fmt.Println("年龄:", user.Age)
   fmt.Println("住址:", user.Address)
   fmt.Println("电话:", user.Phone)
}
```
现在用户信息丰富了，之后还会有批量诉求，要支持传入批量，于是改造如下：
```go
func CollectUserInfoV3(users []*User) {
   for _, user := range users {
      fmt.Println("姓名:", user.Name)
      fmt.Println("年龄:", user.Age)
      fmt.Println("住址:", user.Address)
      fmt.Println("电话:", user.Phone)
   }
}
```
&emsp;&emsp;如上，我们已经拥有了3个版本的方法来支持各种不同差异化的参数了，未来还有扩容的可能性，比如为了检索的便利性，可能需要一个以用户ID为`Key`，用户信息为`Value`的`Map<userId,*UserInfo>`样式的数据结构。

&emsp;至此，我们遇到的问题就是**不能预先确定参数类型，需要动态的执行不同参数类型行为**，这便是反射要做的事情，下面我们引入`Reflect`包来识别不同类型参数，动态的路由到不同的处理单元中解析参数然后执行对应的逻辑，如下：
```go
func CollectUserInfo(param interface{}) {
   val := reflect.ValueOf(param)
   switch val.Kind() {
   case reflect.String:
      fmt.Println("姓名:", val.String())
   case reflect.Struct:
      fmt.Println("姓名:", val.FieldByName("Name"))
      fmt.Println("年龄:", val.FieldByName("Age"))
      fmt.Println("住址:", val.FieldByName("Address"))
      fmt.Println("电话:", val.FieldByName("Phone"))
   case reflect.Ptr:
      fmt.Println("姓名:", val.Elem().FieldByName("Name"))
      fmt.Println("年龄:", val.Elem().FieldByName("Age"))
      fmt.Println("住址:", val.Elem().FieldByName("Address"))
      fmt.Println("电话:", val.Elem().FieldByName("Phone"))
   case reflect.Array, reflect.Slice:
      for i := 0; i < val.Len(); i++ {
         fmt.Println("姓名:", val.Index(i).FieldByName("Name"))
         fmt.Println("年龄:", val.Index(i).FieldByName("Age"))
         fmt.Println("住址:", val.Index(i).FieldByName("Address"))
         fmt.Println("电话:", val.Index(i).FieldByName("Phone"))
      }
   case reflect.Map:
      itr := val.MapRange()
      for itr.Next() {
         fmt.Println("用户ID:", itr.Key())
         fmt.Println("姓名:", itr.Value().Elem().FieldByName("Name"))
         fmt.Println("年龄:", itr.Value().Elem().FieldByName("Age"))
         fmt.Println("住址:", itr.Value().Elem().FieldByName("Address"))
         fmt.Println("电话:", itr.Value().Elem().FieldByName("Phone"))
      }
   default:
      panic("unsupport type !!!")
   }
}

func TestCollectUserInfo(t *testing.T) {
   //string
   CollectUserInfo("张三")
   // 姓名: 张三

   //Struct
   CollectUserInfo(User{
      Name:    "张三",
      Age:     20,
      Address: "北京市海淀区",
      Phone:   1234567,
   })
   //姓名: 张三
   //年龄: 20
   //住址: 北京市海淀区
   //电话: 1234567

   //Ptr
   CollectUserInfo(&User{
      Name:    "张三",
      Age:     20,
      Address: "北京市海淀区",
      Phone:   1234567,
   })
   //姓名: 张三
   //年龄: 20
   //住址: 北京市海淀区
   //电话: 1234567

   //Slice
   CollectUserInfo([]User{
      {
         Name:    "张三",
         Age:     20,
         Address: "北京市海淀区",
         Phone:   1234567,
      },
      {
         Name:    "李四",
         Age:     30,
         Address: "北京市朝阳区",
         Phone:   7654321,
      },
   })
   //姓名: 张三
   //年龄: 20
   //住址: 北京市海淀区
   //电话: 1234567
   //姓名: 李四
   //年龄: 30
   //住址: 北京市朝阳区
   //电话: 7654321

   //Array
   CollectUserInfo([2]User{
      {
         Name:    "张三",
         Age:     20,
         Address: "北京市海淀区",
         Phone:   1234567,
      },
      {
         Name:    "李四",
         Age:     30,
         Address: "北京市朝阳区",
         Phone:   7654321,
      },
   })
   //姓名: 张三
   //年龄: 20
   //住址: 北京市海淀区
   //电话: 1234567
   //姓名: 李四
   //年龄: 30
   //住址: 北京市朝阳区
   //电话: 7654321

   CollectUserInfo(map[int]*User{
      1: {
         Name:    "张三",
         Age:     20,
         Address: "北京市海淀区",
         Phone:   1234567,
      },
      2: {
         Name:    "李四",
         Age:     30,
         Address: "北京市朝阳区",
         Phone:   7654321,
      },
   })
   //用户ID: 1
   //姓名: 张三
   //年龄: 20
   //住址: 北京市海淀区
   //电话: 1234567
   //用户ID: 2
   //姓名: 李四
   //年龄: 30
   //住址: 北京市朝阳区
   //电话: 7654321
}
```
&emsp;&emsp;这样，我们把所有类型的参数都封装到一个函数方法中，未来的变化，扩展化外部是不需要感知的，有没有体会到反射功能的强大呢！

# Reflect源码构成
&emsp;&emsp;**源码之下无秘密**。这里基于`Go1.16`深入源码之中结合注释内容进行层层剖析！`reflect`包下内容可以大体分为三部分：**测试文件**、**编译文件**、**反射核心代码**，这里主要围绕核心代码展开。

<div align=center>
<img src="https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/e907b226b756493e84ad7b68b784be8a~tplv-k3u1fbpfcp-watermark.image?" width="300">
</div>

&emsp;&emsp;按照`reflect`的介绍，这个包主要是实现了**运行时反射**，允许程序操作任意类型的对象。典型用法是：
- 用静态类型`interface{}`保存一个值，通过调用`TypeOf`获取其动态类型信息，该函数返回一个`Type`类型值。
- 调用`ValueOf`函数返回一个`Value`类型值，该值代表运行时的数据。

&emsp;&emsp;也就是说，通过`interface{}`，调用对应方法可以获取到`Type`和`Value`，在`Golang`中`interface{}`可以代表所有类型，那么数据类型与反射之间的关系可以简单总结如下：

<div align=center>
<img src="https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/60ea80bce96e482989f6c9733a769574~tplv-k3u1fbpfcp-watermark.image?" width="400">
</div>

## type.go

<div align=center>
<img src="https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/55b37357f491403b8df33cc2f2196879~tplv-k3u1fbpfcp-watermark.image?" width="700">
</div>

&emsp;&emsp;在**type.go**中最核心内容由`Type`方法定义、`rtype`元数据定义两部分构成。

<div align=center>

定义 | 描述 | 角色 | 可见性 |功能
--|--|--|--|--
Type | 类型能力定义，类型方法的抽象；`Method`、`Kind`、`ChanDir`、`StructField`等属于表达Type能力定义的一部分 | **行为，功能映射** | 可导出，供外部使用 |提供支持所有类型的通用方法、只支持具体类型的方法
rtype | 元数据定义，类型数据的组织结构，衍生出`arrayType`、`chanType`、`mapType`、`sliceType`等常见聚合结构的类型结构定义| **元数据，内存映射** | 不可导出，内部使用 | 作为具体类型方法的基础数据，具体衍生的类型结构也是不可导出的，内部使用

</div>

### Type

`Type`接口定义的源码及注释翻译如下：

```go

// Type类型值是可比较的，例如使用==运算符，因此它们可以用作映射键。
// 如果两个类型值表示相同的类型，则它们相等。
type Type interface {
   // 适用于所有类型的方法.

   // Align在内存中分配时返回此类型值的字节对齐方式
   Align() int

   // FieldAlign在结构中用作字段时，返回此类型值的字节对齐方式
   FieldAlign() int

   // 方法返回类型的方法集中的第i个方法。如果不在范围[0，NumMethod（））内，会导致panic.
   // 对于非接口类型T或*T，返回的方法的type和Func字段描述其第一个参数为接收器的函数，并且只有导出的方法才可访问
   // 对于接口类型，返回方法的类型字段给出了方法签名，没有接收器，而Func字段为零
   // 方法按字典顺序排序
   Method(int) Method

   // MethodByName返回在类型的方法集中具有该名称的方法，并返回一个布尔值，指示是否找到该方法
   // 对于非接口类型T或*T，返回方法的type和Func字段描述其第一个参数为接收器的函数
   // 对于接口类型，返回的方法的类型字段给出了方法签名，没有接收器，而Func字段为零
   MethodByName(string) (Method, bool)

   // NumMethod返回使用方法可访问的方法数
   // 注意，NumMethod仅对接口类型统计未导出的方法
   NumMethod() int

   // Name返回定义类型的包中的类型名称
   // 对于其他（未定义）类型，它返回空字符串
   Name() string

   // PkgPath返回定义类型的包路径，即唯一标识包的导入路径，例如“encoding/base64”
   // 如果类型已预先声明（字符串、错误）或未定义（*T、struct{}、[]int或A，其中A是未定义类型的别名），则包路径将为空字符串
   PkgPath() string

   // Size返回存储给定类型的值所需的字节数；这类似于unsafe.Sizeof()的功能
   Size() uintptr

   // String返回类型的字符串表示形式。
   // 字符串表示可以使用缩短的包名（例如，base64而不是“encoding/base64”），并且不能保证在类型中是唯一的。要测试类型标识，请直接比较类型。
   String() string

   // Kind返回此类型的特定种类
   Kind() Kind

   // 报告类型是否实现接口类型u
   Implements(u Type) bool

   // 报告类型的值是否可分配给类型u
   AssignableTo(u Type) bool

   // 报告类型的值是否可转换为u类型
   ConvertibleTo(u Type) bool

   // 报告此类型的值是否可比较
   Comparable() bool

   // 仅适用于某些类型的Method取决于Kind
   // 每个Kind类型允许的方法如下:
   // Int*, Uint*, Float*, Complex*: Bits
   // Array: Elem, Len
   // Chan: ChanDir, Elem
   // Func: In, NumIn, Out, NumOut, IsVariadic.
   // Map: Key, Elem
   // Ptr: Elem
   // Slice: Elem
   // Struct: Field, FieldByIndex, FieldByName, FieldByNameFunc, NumField

   // Bits返回类型的大小（以位为单位）
   // 如果类型的种类不是大小为或未大小为Int、Uint、Float或Complex的种类之一，则会导致panic。
   Bits() int

   // ChanDir返回通道类型的方向
   // 如果类型不是Chan，它会panic.
   ChanDir() ChanDir

   // IsVariadic报告函数类型的最终输入参数是否为“…”参数如果是，t.In（t.NumIn（）-1）返回参数的隐式实际类型[]T
   // 具体举例，如果t表示func（x int，y…float64），则
   // t.NumIn() == 2
   // t.In(0) is the reflect.Type for "int"
   // t.In(1) is the reflect.Type for "[]float64"
   // t.IsVariadic() == true
 
   // 如果类型的种类不是Func会panic
   IsVariadic() bool

   // Elem返回类型的元素类型
   // 如果类型的种类不是Array、Chan、Map、Ptr或Slice，会出现panic
   Elem() Type

   // 字段返回结构类型的第i个字段
   // 如果类型的种类不是Struct，它就会panic
   // 如果不在范围[0，NumField（））内，它就会panic
   Field(i int) StructField

   // FieldByIndex返回与索引序列对应的嵌套字段。它相当于为每个索引i依次调用字段
   // 如果类型的种类不是Struct，它就会panic
   FieldByIndex(index []int) StructField

   // FieldByName返回具有给定名称和布尔值的结构字段，该布尔值指示是否找到该字段
   FieldByName(name string) (StructField, bool)

   // FieldByNameFunc返回结构字段，其名称满足匹配函数，布尔值指示是否找到该字段
   // FieldByNameFunc首先考虑结构本身中的字段，然后考虑任何嵌入结构中的字段，按广度优先顺序，在包含一个或多个满足匹配函数的字段的最浅嵌套深度处停止
   // 如果该深度的多个字段满足match函数，则它们会相互抵消，FieldByNameFunc不会返回匹配
   // 此行为反映了Go在包含嵌入字段的结构中处理名称查找的方式
   FieldByNameFunc(match func(string) bool) (StructField, bool)

   // In返回函数类型的第i个输入参数的类型
   // 如果类型的种类不是Func，它就会panic
   // 如果不在范围[0，NumIn（））内，它就会panic
   In(i int) Type

   // Key返回映射类型的键类型
   // 如果类型的类型不是Map，它会感到panic
   Key() Type

   // Len返回数组类型的长度
   // 如果类型的种类不是数组，它就会panic
   Len() int

   // NumField返回结构类型的字段计数
   // 如果类型的种类不是Struct，它就会panic
   NumField() int

   // NumIn返回函数类型的输入参数计数
   // 如果类型的种类不是Func，它就会panic
   NumIn() int

   // NumOut返回函数类型的输出参数计数
   // 如果类型的种类不是Func，它就会panic
   NumOut() int

   // Out返回函数类型的第i个输出参数的类型
   // 如果类型的种类不是Func，它就会panic
   // 如果不在范围[0，NumOut（））内，它就会panic
   Out(i int) Type

   common() *rtype
   uncommon() *uncommonType
}
```
在此接口的定义前注释给出了`Type`的一些使用姿势：
- 并非所有方法都适用于所有类型，`Type`定义方法都有其限制条件
- 在调用特定于种类的方法之前，请使用**Kind()** 方法找出类型的种类
- 调用不适合该类型的方法会导致运行时**Panic**

消化下核心要点，`Type`的接口定义并不适合所有类型，那么我们来看下**Kind()** 都定义了哪些类型呢？答案就在`type Kind uint`中，对`Golang`中支持的类型进行了统一定义，如下：

 ```go
type Kind uint

const (
   Invalid Kind = iota
   Bool
   Int
   Int8
   Int16
   Int32
   Int64
   Uint
   Uint8
   Uint16
   Uint32
   Uint64
   Uintptr
   Float32
   Float64
   Complex64
   Complex128
   Array
   Chan
   Func
   Interface
   Map
   Ptr
   Slice
   String
   Struct
   UnsafePointer
)
```
&emsp;&emsp;此外，不同类型`Type`支持不同的方法或行为，需要按照`Kind`列表中具体类型来进行合理调用支持的方法，错误地调用会产生`panic`。

&emsp;&emsp;整理下`Type`给不同类型定义提供了哪些**通用化**或**差异化**的能力定义支持：

<div align=center>

Type定义方法 | 适用类型
--|--
Align() int | 全部
FieldAlign() int | 全部
Method(int) Method | 全部
MethodByName(string) (Method, bool) | 全部
NumMethod() int | 全部
Name() string | 全部
PkgPath() string | 全部
Size() uintptr | 全部
String() string | 全部
Kind() Kind | 全部
Implements(u Type) bool | 全部
AssignableTo(u Type) bool | 全部
ConvertibleTo(u Type) bool | 全部
Comparable() bool | 全部
Bits() int | Int* , Uint* , Float* , Complex*
Elem() Type | Array , Chan , Map , Ptr , Slice
Len() int | Array
ChanDir() ChanDir | Chan
In(i int) Type | Func
NumIn() int | Func
Out(i int) Type |  Func
NumOut() int | Func
IsVariadic() bool | Func
Key() Type | Map
Field(i int) StructField | Struct
FieldByIndex(index []int) StructField | Struct
FieldByName(name string) (StructField, bool) | Struct
FieldByNameFunc(match func(string) bool) (StructField, bool) | Struct
NumField() int | Struct

</div>

### rtype

再来看下`rtype`元数据定义结构：
```go
type rtype struct {
   size       uintptr
   ptrdata    uintptr // rtype可以包含指针的字节数
   hash       uint32  // rtype哈希值；避免哈希表中的计算
   tflag      tflag   // 额外的类型信息标识
   align      uint8   // 当前具体类型变量的内存对齐
   fieldAlign uint8   // 当前具体类型结构体字段的内存对齐
   kind       uint8   // 具体Kind的枚举值
   // 当前具体类型使用的对比方法
   // (ptr to object A, ptr to object B) -> ==?
   equal     func(unsafe.Pointer, unsafe.Pointer) bool
   gcdata    *byte   // 垃圾回收数据
   str       nameOff // 字符串格式
   ptrToThis typeOff // 指向此类型的指针的类型，可以为零
}
```
## value.go
和**type.go**类似，在**value.go**中最核心的是`type Value struct`，它对`Value`进行了抽象定义，这个文件内的其他代码也围绕此来构建做能力支持。`Value`结构体定义的源码及注释翻译如下：

```go
// Value是Go值的反射.

// 并非所有方法都适用于所有类型的值。每种方法的文档中都注明了限制条件（如有）。
// 在调用特定于种类的方法之前，请使用种类方法找出值的种类。调用不适合该类型的方法会导致运行时panic

// 零值代表未赋值、空值
// 零值的IsValid方法返回false，其Kind方法返回Invalid，其String方法返回“<Invalid Value>”，所有其他方法都无法使用
// 大多数函数和方法从不返回无效值
// 如果有，其文档将明确说明这些条件

// 一个值可以由多个goroutine同时使用，前提是基础Go值可以同时用于等效的直接操作

// 要比较两个值，请比较接口方法的结果。在两个值上使用==不会比较它们表示的基础值
type Value struct {
   // typ保存由值表示的值的类型。
   typ *rtype

   // 指针值数据，或者，如果设置了flagIndir，则为指向数据的指针
   // 设置flagIndir或typ.pointers()为true
   // 这是非常核心的数据，可以把它理解为具体数据的内存位置所在，数据的类型表达依赖它来转换
   ptr unsafe.Pointer

   // flag是一个标志位，通过二进制的方式保存了关于值的元数据
   // 最低位是标志位！最低的五位给出了值的类型，代表Kind的枚举的二进制，一共是27个，用5位表示，其余依次如下：
   // - flagStickyRO: 代表不能导出的非嵌入字段获取，因此为只读
   // - flagEmbedRO: 代表不能导出的嵌入字段获取，因此为只读
   // - flagIndir: 代表持有指向数据的指针
   // - flagAddr: 代表CanAddr方法的返回值标记
   // - flagMethod: 代表是否为一个方法的标记
   // 剩余的高23位给出了方法值的方法编号。
   // 如果flag.Kind（）！=Func，代码可以假设flagMethod未设置。
   // 如果ifaceIndir（typ）为真，则代码可以假设设置了flagIndir。
   flag

   // 方法值表示当前的方法调用，就像接收者r调用r.Read方法。typ+val+flag的标志位描述r的话，但flag的Kind标志位表示Func（方法是函数），flag的高位表示r类型的方法列表中的方法编号
}
```

可以看到，**Value**主要由`typ(*rtype)`、`ptr（unsafe.Pointer）`、`flag(uintptr)`构成。
<div align=center>

组成 | 功能
--|--
typ(*rtype) | 数据存储，内存映射
ptr（unsafe.Pointer） | 指针
flag(uintptr) | 二进制标记

</div>

来看下`flag`的枚举定义以及标志位的二进制占位分布情况：

<div align=center>
<img src="https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/46101cbe437c4299ab3784a3afd7cf58~tplv-k3u1fbpfcp-watermark.image?" width="600">
</div>

```go
type flag uintptr

const (
   flagKindWidth        = 5 // 有27个Kind类型，5位可以容纳2^5可以表示32个类型
   flagKindMask    flag = 1<<flagKindWidth - 1
   flagStickyRO    flag = 1 << 5
   flagEmbedRO     flag = 1 << 6
   flagIndir       flag = 1 << 7
   flagAddr        flag = 1 << 8
   flagMethod      flag = 1 << 9
   flagMethodShift      = 10
   flagRO          flag = flagStickyRO | flagEmbedRO
)
```
类似地，**Value**也提供了**通用化**或**差异化**的能力定义支持：

Value定义方法 | 适用类型
--|--
Addr() Value |全部
Bool() bool | Bool
Bytes() []byte |  Slice
CanAddr() bool |  全部
CanSet() bool |  全部
Call(in []Value) []Value |  Func
CallSlice(in []Value) []Value |  Func
Cap() int |  Array , Chan , Slice
Close() |  Chan
Complex() complex128 |  Complex64 , Complex128
Elem() Value |  Interface , Ptr
Field(i int) Value |  Struct
FieldByIndex(index []int) Value |  Struct
FieldByName(name string) Value |  Struct
FieldByNameFunc(match func(string) bool) Value |  Struct
Float() float64 |  Float32 , Float64
Index(i int) Value |  Array , Slice , String
Int() int64 |  Int*
CanInterface() bool |  全部
Interface() (i interface{}) |  Interface
InterfaceData() [2]uintptr |  Interface
IsNil() bool |  Chan , Func , Map , Ptr , UnsafePointer , Interface , Slice
IsValid() bool |  全部
IsZero() bool |  全部
Kind() Kind |  全部
Len() int |  Array , Chan , Map , Slice , String
MapIndex(key Value) Value |  Map
MapKeys() []Value |  Map
MapRange() *MapIter |  Map
Method(i int) Value |  全部
NumMethod() int |  全部
MethodByName(name string) Value |  全部
NumField() int |  Struct
OverflowComplex(x complex128) bool |  Complex64 , Complex128
OverflowFloat(x float64) bool |  Float32 , Float64
OverflowInt(x int64) bool |  Int*
OverflowUint(x uint64) bool |  Uint*
Pointer() uintptr |  Ptr , Chan , Map , UnsafePointer , Func , Slice
Recv() (x Value, ok bool) |  Chan
Send(x Value) |  Chan
Set(x Value) |  全部
SetBool(x bool) | Bool
SetBytes(x []byte) | Slice
SetComplex(x complex128) |  Complex64 , Complex128
SetFloat(x float64) | Float32 , Float64
SetInt(x int64) | Int*
SetLen(n int) | Slice
SetCap(n int) | Slice
SetMapIndex(key Value, elem Value) | Map
SetUint(x uint64) | Uint*
SetPointer(x unsafe.Pointer) | UnsafePointer
SetString(x string) | String
Slice(i int, j int) Value | Array , Slice , String
Slice3(i int, j int, k int) Value |  Array , Slice , String
String() string | String , Invalid
TryRecv() (x Value, ok bool) | Chan
TrySend(x Value) bool | Chan
Type() Type | 全部
Uint() uint64 | Uint*
UnsafeAddr() uintptr | 全部
Convert(t Type) Value | 全部

## makefunc.go
`makefunc`提供了函数反射调用的能力：

```go
// MakeFunc返回一个给定类型的新函数，该函数封装了函数fn。调用时，该新函数执行以下操作：
// - 将其参数转换为一个切片Slice的值。
// - 运行结果：=fn（args）。
// - 将结果作为一个切片Slice的值返回，每个正式结果一个值。

// 实现fn可以假设参数值切片具有typ给定的参数数量和类型。
// 如果typ描述可变函数，则最终值本身是表示可变参数的切片，如在可变函数体中。
// fn返回的结果值切片必须具有typ给定的结果数量和类型。

// Value.Call方法允许调用方根据值调用类型化函数；相反，MakeFunc允许调用方根据值实现类型化函数。

// 文档的示例部分演示了如何使用MakeFunc为不同类型构建交换函数
func MakeFunc(typ Type, fn func(args []Value) (results []Value)) Value {
   if typ.Kind() != Func {
      panic("reflect: call of MakeFunc with non-Func type")
   }

   t := typ.common()
   ftyp := (*funcType)(unsafe.Pointer(t))

   // Go func的间接值（虚拟）以获取实际代码地址。
   // Go func值是指向C函数指针的指针。https://golang.org/s/go11func
   dummy := makeFuncStub
   code := **(**uintptr)(unsafe.Pointer(&dummy))

   // makeFuncImpl包含一个堆栈映射，供运行时使用
   _, argLen, _, stack, _ := funcLayout(ftyp, nil)

   impl := &makeFuncImpl{code: code, stack: stack, argLen: argLen, ftyp: ftyp, fn: fn}

   return Value{t, unsafe.Pointer(impl), flag(Func)}
}
```

## swapper.go
&emsp;&emsp;标准库实现了一个支持`Slice`切片元素按照索引进行交换的方法，它的底层完全是由`Reflect`包能力来支持的，可以把它作为一个`Reflect`使用范例来学习！

```go
func TestSwapper(t *testing.T) {
   s := []int{1,2,3,4,5} // 声明一个切片，元素排列为 [1 2 3 4 5]
   f := reflect.Swapper(s) // 调用reflect.Swapper()方法，出参是一个方法
   f(0,1)          // 调用方法，将索引位 0、1的元素互换
   fmt.Println(s) // 结果为[2 1 3 4 5]
}
```

## deepequal.go
&emsp;&emsp;标准库还实现了一个支持深度比较相等的方法，它的底层也完全是由`Reflect`包能力来支持的，同样可以把它作为一个`Reflect`使用范例来学习！

```go
func TestDeepEqual(t *testing.T) {
   x := &User{Name: "zhangsan", Age: 10}
   y := &User{Name: "zhangsan", Age: 10}
   fmt.Println(reflect.DeepEqual(x, y)) // true

   x1 := &User{Name: "zhangsan", Age: 10}
   y1 := &User{Name: "zhangsan", Age: 0}
   fmt.Println(reflect.DeepEqual(x1, y1)) // false

   x2 := map[string]int{"zhangsan": 10}
   y2 := map[string]int{"zhangsan": 10}
   fmt.Println(reflect.DeepEqual(x2, y2)) // true

   x3 := map[string]int{"zhangsan": 10}
   y3 := map[string]int{"lisi": 10}
   fmt.Println(reflect.DeepEqual(x3, y3)) // false
}
```

# 反射三大定律
> *源码中提到了**Rob Pike**写的关于`Go`中反射介绍的文章，里面提到使用反射机制的一些规则，可以作为范本规约来从更高视角理解反射。 [The Laws of Reflection](https://go.dev/blog/laws-of-reflection)*
## Reflection goes from interface value to reflection object
> *反射可以将 **“接口类型变量”** 转换为 **“反射类型对象”***

```go
func TypeOf(i interface{}) Type

func ValueOf(i interface{}) Value
```

```go
func TestInterface2ReflectionObj(t *testing.T) {
   f := float32(3.14)
   // Interface 转换为 Type
   typ := reflect.TypeOf(f)
   // Interface 转换为 Value
   val := reflect.ValueOf(f)
   fmt.Println(typ) // float32
   fmt.Println(val) // 3.14
}
```

&emsp;&emsp;**反射**只是一种检查存储在接口变量中的类型和值对的机制。首先，我们需要了解[反射包](https://go.dev/pkg/reflect/)中的两种类型： [类型](https://go.dev/pkg/reflect/#Type)和[值](https://go.dev/pkg/reflect/#Value)。这两种类型可以访问接口变量的内容，以及两个简单的函数，分别是`reflect.TypeOf`and `reflect.ValueOf`。

## Reflection goes from reflection object to interface value
> *反射可以将 **“反射类型对象”** 转换为 **“接口类型变量”***

```go
func (v Value) Interface() (i interface{})
```

```go
func TestReflectionObj2Interface(t *testing.T) {
   f := float32(3.14)
   // 通过Reflect对象转换Interface
   val := reflect.ValueOf(f)
   // 转换指定的Interface
   fmt.Printf("%T %f \n",val.Interface().(float32),val.Interface().(float32)) // float32 3.140000
}
```
## To modify a reflection object， the value must be settable
> *想要修改反射对象，它的值必须是可赋值的（可写的）*

&emsp;&emsp;什么是可赋值？它依赖`Value.CanSet()`方法提供判断，其内部是根据`flagAddr`进行判断的，这里参考[反射定律](https://go.dev/blog/laws-of-reflection)的描述进行了一些可赋值的测试代码：

```go
func TestSettable(t *testing.T) {
   //****可赋值****
   //声明一个对象，并初始化赋值
   user := &User{Name: "zhangsan",Age: 10}
   // 判断是否可赋值
   fmt.Println(reflect.ValueOf(user).CanSet())        // false
   fmt.Println(reflect.ValueOf(user).Elem().CanSet()) // true
   // 获取对象user中Name字段的值
   fmt.Println(reflect.ValueOf(user).Elem().FieldByName("Name")) // zhangsan
   // 获取对象user中Age字段的值
   fmt.Println(reflect.ValueOf(user).Elem().FieldByName("Age")) // 10
   reflect.ValueOf(user).Elem().Field(0).SetString("lisi")
   reflect.ValueOf(user).Elem().Field(1).SetInt(20)
   fmt.Println(user) // &{lisi 20}

   //****不可赋值****
   f := float32(3.14)
   // 这里的f是一个临时变量声明，由于安全性和变更的未知性带来的潜在问题，这里是不可寻址的
   fmt.Println(reflect.ValueOf(f), reflect.ValueOf(f).CanSet()) // 3.14 false
   // &f是指向变量f的指针变量，对于此和f的变量的不可寻址本质是一样的
   fmt.Println(reflect.ValueOf(&f), reflect.ValueOf(&f).CanSet()) // 0xc0000182e8 false
   //****可赋值****
   //这里调用Elem主要是复制了&f地址空间，并使返回的Elem(&f)变为可赋值！
   fmt.Println(reflect.ValueOf(&f).Elem(), reflect.ValueOf(&f).Elem().CanSet()) // 3.14 true
}
```
&emsp;&emsp;不难发现，都是通过调用`Value.Elem()`进行了可赋值操作，该方法内部重新拷贝了原`Value`结构，并给标志位`flagAddr`打标，使其变为**可赋值**。这样做的好处是可以通过该属性`flagAddr`来灵活控制是否可赋值，而不是作为一种默认功能来支持，避免在变量赋值、拷贝等处理过程中产生歧义和问题，具体可以参考[反射定律](https://go.dev/blog/laws-of-reflection)中对该问题的阐述。

# 反射的应用
&emsp;&emsp;反射在Golang的应用非常广泛，这里展示一段标准库中`encoding/json`的片段，大量使用到`Reflect`包来进行类型映射和动态处理
```go
// newTypeEncoder constructs an encoderFunc for a type.
// The returned encoder only checks CanAddr when allowAddr is true.
func newTypeEncoder(t reflect.Type, allowAddr bool) encoderFunc {
   // If we have a non-pointer value whose type implements
   // Marshaler with a value receiver, then we're better off taking
   // the address of the value - otherwise we end up with an
   // allocation as we cast the value to an interface.
   if t.Kind() != reflect.Ptr && allowAddr && reflect.PtrTo(t).Implements(marshalerType) {
      return newCondAddrEncoder(addrMarshalerEncoder, newTypeEncoder(t, false))
   }
   if t.Implements(marshalerType) {
      return marshalerEncoder
   }
   if t.Kind() != reflect.Ptr && allowAddr && reflect.PtrTo(t).Implements(textMarshalerType) {
      return newCondAddrEncoder(addrTextMarshalerEncoder, newTypeEncoder(t, false))
   }
   if t.Implements(textMarshalerType) {
      return textMarshalerEncoder
   }

   switch t.Kind() {
   case reflect.Bool:
      return boolEncoder
   case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
      return intEncoder
   case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
      return uintEncoder
   case reflect.Float32:
      return float32Encoder
   case reflect.Float64:
      return float64Encoder
   case reflect.String:
      return stringEncoder
   case reflect.Interface:
      return interfaceEncoder
   case reflect.Struct:
      return newStructEncoder(t)
   case reflect.Map:
      return newMapEncoder(t)
   case reflect.Slice:
      return newSliceEncoder(t)
   case reflect.Array:
      return newArrayEncoder(t)
   case reflect.Ptr:
      return newPtrEncoder(t)
   default:
      return unsupportedTypeEncoder
   }
}
```
# 总结
<div align=center>
<img src="https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/41e2c8a5eae94a52bddaaafc98bf39b9~tplv-k3u1fbpfcp-watermark.image?" width="500">
</div>

&emsp;&emsp;可以从三层视角来理解反射在`Golang`中发挥的作用：
- **业务代码**

  &emsp;&emsp;业务代码的数据组织方式完全交由开发者决定，基于`Golang`提供的数据类型构建出丰富的业务功能，这是现实业务与编程语言的抽象映射。

- **标准库**

  &emsp;&emsp;`Golang`中数据类型是多样、复杂的，它通过定义通用的数据结构和组织方式将复杂类型从朴素的内存空间的存储单元中抽象、拼装出来，具体存储的位置和数据类型是无关的，只关乎于数据是否匹配预先设置的格式空间要求，一旦满足即可完成内存到丰富数据表现的转换能力！

  &emsp;&emsp;`interface{}`更面向对象角度便于开发者使用和操作，更偏向于上层，无需关心`字符串、数值、布尔`这些数据类型是如何存储和内存转换的，直接操作这些`”成品“`即可；而`reflect`更面向数据存储、内存映射便于裸数据到具体上游数据类型的转换，更偏向于底层，这些数据组织方式更像是`半成品`，还有很多可塑和操作空间。`reflect`提供了`TypeOf、ValueOf`方法与`interface{}`进行转换和数据交互。

- **内存空间**

  &emsp;&emsp;内存空间是程序在运行时客观驻留的物理空间，无论在Golang中是`字符串、数值、布尔`这种简单的数据类型，还是`Map、结构体、数组、切片`等复杂的数据类型，在内存空间的表现都可以统一表示为`偏移量`

## 优点
- 避免硬编码，提供灵活性和通用性
- 可以作为第一类对象发现并修改源码结构（代码块、类、方法、协议等）

## 缺点
- **可读性差** 反射不同于一般的业务代码那样容易理解和通用，需要一定的背景知识才可以读懂它
- **错误风险** 在Go语言中，`interface`的定义类型是不确定的，带来便利性的同时也伴随着`panic`风险
- **性能开销** 反射需要进行动态地类型匹配和寻址操作，存在一定性能损耗

# 参考
[Golang标准库文档](https://studygolang.com/pkgdoc)\
[Go夜谈 #108 Golang 反射应用及源码分析](https://www.bilibili.com/video/BV1My4y117gQ)\
[Go官方反射介绍](https://go.dev/blog/laws-of-reflection)
《Go语言圣经》反射章节 \
[手摸手Go 接口与反射](https://www.jianshu.com/p/6c67741a5c52)