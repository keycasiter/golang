package code

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestCopyFunc(t *testing.T) {
	obj := make([]*Element, 0)
	copyObj := make([]*Element, 1)

	obj = append(obj, &Element{
		Name: "111",
		Element: &Element{
			Name:    "222",
			Element: nil,
		},
	})

	fmt.Println(copy(copyObj, obj))
	fmt.Println(unsafe.Pointer(obj[0]), unsafe.Pointer(copyObj[0]))
	fmt.Println("=================")
	fmt.Println(unsafe.Pointer(obj[0].Element), unsafe.Pointer(copyObj[0].Element))
}
