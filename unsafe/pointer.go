package main

import (
	"fmt"
	"unsafe"
)

/**
 * unsafe.Pointer 是一个通用性的指针, 可以在任何类型的指针之间转换
 * 关于unsafe.Pointer的四个规则:
 * 1. 任何指针都可以转换为unsafe.Pointer
 * 2. unsafe.Pointer可以转换为任何指针
 * 3. uintptr可以转换为unsafe.Pointer
 * 4. unsafe.Pointer可以转换为uintptr
 *
 * 前面两个规则我们刚刚已经演示了,主要用于*T1和*T2之间的转换,那么最后两个规则是做什么的呢？
 * 我们都知道*T是不能计算偏移量的,也不能进行计算,但是uintptr可以.
 * 所以我们可以把指针转为uintptr再进行偏移计算,这样我们就可以访问特定的内存了,达到对不同的内存读写的目的
 */

// 下面我们以通过指针偏移修改Struct结构体内的字段为例，来演示uintptr的用法。
func main() {
	u := new(user)
	fmt.Println(*u)

	pName := (*string)(unsafe.Pointer(u))
	*pName = "张三"

	pAge := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(u)) + unsafe.Offsetof(u.age)))
	*pAge = 20

	fmt.Println(*u)
}

type user struct {
	name string
	age  int
}

/**
以上我们通过内存偏移的方式，定位到我们需要操作的字段，然后改变他们的值。

第一个修改user的name值的时候，因为name是第一个字段，所以不用偏移，我们获取user的指针，然后通过unsafe.Pointer转为*string进行赋值操作即可。

第二个修改user的age值的时候，因为age不是第一个字段，所以我们需要内存偏移，内存偏移牵涉到的计算只能通过uintptr，所我们要先把user的指针地址转为uintptr，然后我们再通过unsafe.Offsetof(u.age)获取需要偏移的值，进行地址运算(+)偏移即可。

现在偏移后，地址已经是user的age字段了，如果要给它赋值，我们需要把uintptr转为*int才可以。所以我们通过把uintptr转为unsafe.Pointer,再转为*int就可以操作了。
*/
