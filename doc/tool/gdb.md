# []()前言

学会使用`gdb`进行`golang`的调试，通过一个简单的go程序进行简单入口程序的源码调用顺序的查看。

# []()gdb安装

开发环境是Mac，可以使用`brew`来进行`gdb`安装

## []()更新brew

```bash
brew update
```

## []()查看是否存在gdb镜像

```bash
brew search gdb
```

## []()安装gdb

```bash
brew install gdb
```

# []()go build编译

关于go的编译命令帮助可以通过`go help build`查看命令参数详情。

这里使用`go build -gcflags=all="-N -l" -ldflags=-compressdwarf=false` 进行编译，这里编译参数含义是关闭函数内联和编译器的代码优化，使得编译可以不受其干扰更好地可以查看编译源码的展示。

这里编译用到的参数含义解释如下：

| 编译参数                          | 含义               |
| ----------------------------- | ---------------- |
| -gcflags=all="-N -l"          | 关闭内联优化；关闭编译器函数优化 |
| -ldflags=-compressdwarf=false | 生成压缩可调式信息           |

# []()gdb执行

执行`gdb main`进入gdb环境，即可进行`gdb`环境的调试和使用

```bash
$ gdb program
Reading symbols from program...
Loading Go Runtime support.
```

# []()gdb命令

在进入`gdb`调试之前，先科普下常用的`gdb`调试命令，如下

| 命令         | 说明                        |
| ---------- | :------------------------ |
| file <文件名> | 加载被调试的可执行程序文件             |
| run        | 重新开始运行文件，简写r              |
| start      | 单步执行，运行程序，停在第一执行语句        |
| list       | 查看源代码，简写l                 |
| set        | 设置变量的值                    |
| next       | 单步调试（逐过程，函数直接执行）,简写n      |
| step       | 单步调试（逐语句：跳入自定义函数内部执行）,简写s |
| backtrace  | 查看函数的调用的栈帧和层级关系,简写bt      |
| frame      | 切换函数的栈帧,简写f               |
| info       | 查看函数内部局部变量的数值,简写i         |
| finish     | 结束当前函数，返回到函数调用点           |
| continue   | 继续运行,简写c                  |
| print      | 打印值及地址,简写p                |
| quit       | 退出gdb,简写q                 |

以上命令的跟随参数可以在`gdb`环境中执行`help [command]`来查看

# []()gdb调试

下面，基于一个简单的golang程序进行gdb调试，目的是了解golang程序的初始化流程。

执行`list`查看当前文件的源代码，这里看到的是我们编写的原始的`go`程序文件

```bash
(gdb) list
1       package main
2       
3       import "fmt"
4       
5       func main() {
6               fmt.Println("hello world")
7       }
```

执行`info files`查看函数内部局部变量的数值，如下

```bash
(gdb) info files
Symbols from "/Users/guanjian/workspace/go/program/program".
Local exec file:
        `/Users/guanjian/workspace/go/program/program', file type mach-o-x86-64.
        Entry point: 0x10701e0:
        0x0000000001001000 - 0x00000000010cbc53 is .text
        0x00000000010cbc60 - 0x00000000010cbd62 is __TEXT.__symbol_stub1
        0x00000000010cbd80 - 0x0000000001100b19 is __TEXT.__rodata
        0x0000000001100b20 - 0x00000000011012bc is __TEXT.__typelink
        0x00000000011012c0 - 0x0000000001101328 is __TEXT.__itablink
        0x0000000001101328 - 0x0000000001101328 is __TEXT.__gosymtab
        0x0000000001101340 - 0x000000000115dae8 is __TEXT.__gopclntab
        0x000000000115e000 - 0x000000000115e020 is __DATA.__go_buildinfo
        0x000000000115e020 - 0x000000000115e178 is __DATA.__nl_symbol_ptr
        0x000000000115e180 - 0x000000000116c4a4 is __DATA.__noptrdata
        0x000000000116c4c0 - 0x00000000011738f0 is .data
        0x0000000001173900 - 0x00000000011a1170 is .bss
        0x00000000011a1180 - 0x00000000011a62f0 is __DATA.__noptrbss
(gdb) 
```

这里看到最上面的入口点是`Entry point: 0x10701e0:`，我们在这个入口进行断点设置进行这个程序的调试，执行`break *0x10701e0:`（这里是根据数据地址进行端点，前面加了星号 \*）

```bash
(gdb) break *0x10701e0
Breakpoint 1 at 0x10701e0: file /usr/local/go/src/runtime/rt0_darwin_amd64.s, line 8.
```

可以看到已经添加断点成功了，然后看到了当前编译的go程序入口在`/usr/local/go/src/runtime/rt0_darwin_amd64.s`文件，我们进入到该目录下查看所有`rt0`开头的文件全是`.s`结尾的汇编语言实现的，如下：

```bash
$ ls /usr/local/go/src/runtime/rt0_
rt0_aix_ppc64.s        rt0_ios_arm64.s        rt0_netbsd_arm.s
rt0_android_386.s      rt0_js_wasm.s          rt0_netbsd_arm64.s
rt0_android_amd64.s    rt0_linux_386.s        rt0_openbsd_386.s
rt0_android_arm.s      rt0_linux_amd64.s      rt0_openbsd_amd64.s
rt0_android_arm64.s    rt0_linux_arm.s        rt0_openbsd_arm.s
rt0_darwin_amd64.s     rt0_linux_arm64.s      rt0_openbsd_arm64.s
rt0_darwin_arm64.s     rt0_linux_mips64x.s    rt0_openbsd_mips64.s
rt0_dragonfly_amd64.s  rt0_linux_mipsx.s      rt0_plan9_386.s
rt0_freebsd_386.s      rt0_linux_ppc64.s      rt0_plan9_amd64.s
rt0_freebsd_amd64.s    rt0_linux_ppc64le.s    rt0_plan9_arm.s
rt0_freebsd_arm.s      rt0_linux_riscv64.s    rt0_solaris_amd64.s
rt0_freebsd_arm64.s    rt0_linux_s390x.s      rt0_windows_386.s
rt0_illumos_amd64.s    rt0_netbsd_386.s       rt0_windows_amd64.s
rt0_ios_amd64.s        rt0_netbsd_amd64.s     rt0_windows_arm.s
```

接下来我们查看`rt0_darwin_amd64.s`文件，如下：

```bash
% cat -n rt0_darwin_amd64.s
     1	// Copyright 2009 The Go Authors. All rights reserved.
     2	// Use of this source code is governed by a BSD-style
     3	// license that can be found in the LICENSE file.
     4
     5	#include "textflag.h"
     6
     7	TEXT _rt0_amd64_darwin(SB),NOSPLIT,$-8
     8		JMP	_rt0_amd64(SB)	//这里上一步的断点位置
     9
    10	// When linking with -shared, this symbol is called when the shared library
    11	// is loaded.
    12	TEXT _rt0_amd64_darwin_lib(SB),NOSPLIT,$0
    13		JMP	_rt0_amd64_lib(SB)
```

第8行的`_rt0_amd64(SB)` 这里是调用入口，即

```bash
Breakpoint 1 at 0x10701e0: file /usr/local/go/src/runtime/rt0_darwin_amd64.s, line 8.
```

下面我们需要查找`_rt0_amd64`这个汇编方法的位置了，查找方法可以根据自己的方式来查看，我这里为了更方便查看，下载了golang的sdk，然后导入IDE中进行全局搜索，定位到文件`runtime/asm_amd64.s`，如下：\
![在这里插入图片描述](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/cdbe546bbc764cd6b65b350d10991021~tplv-k3u1fbpfcp-zoom-1.image)\
接着，我们再来查找`runtime·rt0_go(SB)`方法，如下\
![在这里插入图片描述](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/56903ae6db804ec1a8ae2089c2ba13e4~tplv-k3u1fbpfcp-zoom-1.image)\
`runtime·rt0_go(SB)`方法包含了很多初始化方法，这里摘下雨痕文章中的示例如下：

```bash
TEXT runtime·rt0_go(SB),NOSPLIT,$0
...
// ===== 调用初始化函数 =====
CALL runtime·args(SB)
CALL runtime·osinit(SB)
CALL runtime·schedinit(SB)
// ===== 创建 main goroutine 用于执行 runtime.main =====
MOVQ $runtime·mainPC(SB), AX
PUSHQ AX
PUSHQ $0
CALL runtime·newproc(SB)
POPQ AX
POPQ AX
// ===== 让当前线程开始执行 main goroutine =====
CALL runtime·mstart(SB)
RET
DATA runtime·mainPC+0(SB)/8,$runtime·main(SB)
GLOBL runtime·mainPC(SB),RODATA,$8
```

到这里，以上所有的go程序入口调用都是在汇编文件下进行的，后面将进入go语言的世界，通过go文件来编写，`$runtime·main(SB)`是函数出口，由此进入go语言的编写逻辑。我们通过gdb断点来查看调用，如下：

```bash
(gdb) b runtime.main
Breakpoint 1 at 0x103b4c0: file /usr/local/go/src/runtime/proc.go, line 115.
```

至此，我们已经搭建了`gdb`的调试环境和一些查看golang函数代码的实践。

# []()问题整理

* **问题1: macOs下gdb签名限制**

```bash
(gdb) run
Starting program: /xxx/main 
Unable to find Mach task port for process-id 35564: (os/kern) failure (0x5).
(please check gdb is codesigned - see taskgated(8))
```

简单的方法可以通过`root`来解决签名受限问题，但是这样调试程序风险很高，可以参考[mac下如何对gdb进行证书签名](https://jingyan.baidu.com/article/c85b7a6437ee0d403bac95b2.html)

* **问题2: (No debugging symbols found in ./main) 而且`(gdb) list` 找不到调试文件**

```bash
(gdb) list
No symbol table is loaded.  Use the "file" command.
```

我这边的解决方法是在编译go时，增加参数`go build -ldflags=-compressdwarf=false`可以输出调试信息便可以了

* **问题3: 关于函数调用查找方法**

1. IDE导入go的sdk进行全局搜索（比较方便，但是不一定准确，需要充分理解调用链路）
2. 进入gdb环境，通过`  b [file].[method] `方式进行断点标记，可以显示当前断点的源文件位置包含行号，也可以通过`info breakpoints`查看断点信息

* **问题4: gdb断点的行号显示在IDE中不准确**\
  可以打开系统中的命令行进行查看，一般是准确的，使用`cat -n file | head -n [number]` 命令即可

* **问题5: gdb 执行run后卡住 & 打断点不生效的可能原因**
  安装完gdb后应该进行代码签名，这里是个gdb在系统中的权限问题，参考[gdb 初次运行卡住 Starting program: [New Thread 0x1103 of process 843]](https://blog.csdn.net/LU_ZHAO/article/details/104810246)

# []()参考

《Go语言学习笔记》雨痕\
[go build 命令参数详解](https://blog.csdn.net/cyq6239075/article/details/103911098)\
[使用GDB调试Go语言](https://www.cnblogs.com/beyondblog/p/4423173.html)\
[mac下如何对gdb进行证书签名](https://jingyan.baidu.com/article/c85b7a6437ee0d403bac95b2.html)\
[Go使用gdb调试](https://jiajunhuang.com/articles/2020\_04\_23-go_gdb.md.html)\
[Debugging Go Code with GDB](https://golang.org/doc/gdb)\
[使用 GDB 调试 Go 代码](https://golang.org/doc/gdb)\
[【Go源码阅读】使用GDB调试Go程序](https://www.bilibili.com/video/BV1rT4y1g77P?from=search\&seid=3745091901446427851)\
[gdb命令](https://linux265.com/course/linux-command-gdb.html)\
[使用gdb添加断点的几种方式](https://www.cnblogs.com/xl2432/p/11361997.html)
[gdb 初次运行卡住 Starting program: [New Thread 0x1103 of process 843]](https://blog.csdn.net/LU_ZHAO/article/details/104810246)
