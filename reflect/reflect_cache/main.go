package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type User struct {
	Name string
	Age  int
}

var handler = func(u *User, message string) {
	fmt.Printf("Hello, My name is %s, I am %d years old ! so, %s\n", u.Name, u.Age, message)
}

//使用普通反射的方式处理名字屏蔽
func filtName(u *User, message string) {
	fn := reflect.ValueOf(handler)
	uv := reflect.ValueOf(u)
	name := uv.Elem().FieldByName("Name")
	name.SetString("XXX")
	fn.Call([]reflect.Value{uv, reflect.ValueOf(message)})
}

//重用部分数据减少重复创建的反射方式处理名字屏蔽
var offset uintptr = 0xFFFF

func filtNameWithReuseOffset(u *User, message string) {
	if offset == 0xFFFF {
		t := reflect.TypeOf(u).Elem()
		name, _ := t.FieldByName("Name")
		offset = name.Offset

	}
	up := (*[2]uintptr)(unsafe.Pointer(&u))
	upnamePtr := (*string)(unsafe.Pointer(up[0] + offset))
	*upnamePtr = "YYY"
	fn := reflect.ValueOf(handler)
	uv := reflect.ValueOf(u)
	fn.Call([]reflect.Value{uv, reflect.ValueOf(message)})
}

func main() {
	u1 := &User{
		Name: "solo",
		Age:  11,
	}
	filtName(u1, "test1")

	u2 := &User{
		Name: "solo",
		Age:  11,
	}
	filtNameWithReuseOffset(u2, "test2")
}
