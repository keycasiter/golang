package code

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"testing"

	"golang.org/x/tools/go/ast/astutil"
)

//打印遍历AST结构树
func TestAst(t *testing.T) {
	var srcCode = `
			package hello
			
			import "fmt"
		
			func greet() {
				var msg = "Hello World!"
				fmt.Println(msg)
			}
			`

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", srcCode, parser.ParseComments)
	if err != nil {
		fmt.Printf("err = %s", err)
	}
	//*** 打印AST结构树 ***
	//ast.Print(fset, f)

	//*** 遍历AST结构树 ***
	// Inspect the AST and print all identifiers and literals.
	ast.Inspect(f, func(n ast.Node) bool {
		var s string
		switch x := n.(type) {
		case *ast.BasicLit:
			s = x.Value
		case *ast.Ident: // 获取属性
			s = x.Name // 此处获取
		}
		if s != "" {
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), s)
		}
		return true
	})
}

//利用astutil修改代码结构
func TestRewrite(t *testing.T) {
	var src string = `package main

import(
	"errors"
)

func main(){
	errors.New("xxxxx")
	errors.New("-----")
	errors.New(Foo("asdfasdf"))
}

func Ts(){
	errors.New("zzz")
}

func Foo(s string){}
`

	// 设置可写入文件
	fset := token.NewFileSet()
	// 解析文件
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		t.Fail()
	}

	//ast.Print(fset, file)
	astutil.RewriteImport(fset, file, "errors", "syserror")

	// 导入file，最终返回file。 .Apply方法类似ast.Inspect
	astutil.Apply(file, nil, func(c *astutil.Cursor) bool {

		n := c.Node()
		switch t := n.(type) {
		case *ast.CallExpr:
			if sl, ok := t.Fun.(*ast.SelectorExpr); ok {
				// 检查语句是否是errors.New
				if sl.Sel.Name != "New" {
					break
				}
				// 断言是否有子属性 sl.X.(*ast.Ident)
				if cl, ok := sl.X.(*ast.Ident); ok {
					// 判断是否是errors，否则跳过
					if cl.Name != "errors" {
						break
					}
					// 替换成syserror
					cl.Name = "syserror"
					// 拼接两个数组
					t.Args = append([]ast.Expr{&ast.BasicLit{Kind: token.IDENT, Value: "traceID"}}, t.Args...)
				}
			}
		}
		return true
	})

	if err := format.Node(os.Stdout, token.NewFileSet(), file); err != nil {
		log.Fatalln("Error:", err)
	}
}
