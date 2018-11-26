package main

import (
	"fmt"
	"runtime"
	"sync"
)

type Msg interface {
	Code() int
	Data() []byte
}

type TextMsg struct {
	content string
}

func (text *TextMsg) Data() []byte {
	return []byte(text.content)
}

func (text *TextMsg) Code() int {
	return len(text.content)
}

func printMsgCode(msg Msg) {
	fmt.Println(msg.Code())
}

//type MyType []int
//
//func add(myType MyType) MyType {
//	for i, _ := range myType {
//		myType[i] += 100
//	}
//	return myType
//}

type show interface {
	showme()
}

// user在程序里定义一个用户类型
type user struct {
	name  string
	email string
}

func (u *user) showme() {
	fmt.Printf("name %s, email %s\n", u.name, u.email)
}

// admin代表一个拥有权限的管理员用户
type admin struct {
	user // 嵌入类型
	level int
}

type student struct {
	u     user
	level int
}


var wg sync.WaitGroup

// printPrime 显示5000以内的素数值
func printPrime(prefix string) {
	// 在函数退出时调用Done来通知main函数工作已经完成
	defer wg.Done()
next:
	for outer := 2000000; outer < 2001000; outer++ {
		for inner := 2; inner < outer; inner++ {
			if outer%inner == 0 {
				continue next
			}
		}
		fmt.Printf("%s:%d\n", prefix, outer)
	}
	fmt.Println("Completed", prefix)
}

func main() {

	runtime.GOMAXPROCS(1)
	wg.Add(2)

	fmt.Println("create goroutines")

	go printPrime("##")
	go printPrime("..")

	fmt.Println("waiting to finish")

	wg.Wait()

	fmt.Println("Done")

	//ad := admin{
	//	user: user{
	//		name:  "Admin",
	//		email: "www@123.com",
	//	},
	//	level: 10,
	//}
	//
	//ad.user.showme()
	//ad.showme()
	//
	//s := student{
	//	u: user{
	//		name:  "student",
	//		email: "student@123.com",
	//	},
	//	level: 10,
	//}
	//
	//s.u.showme()

	//msg := TextMsg{content: "str"}
	//
	//fmt.Println(msg.Code())
	//fmt.Println(msg.Data())
	//
	//printMsgCode(&msg)

	//var myType MyType = make(MyType, 10)
	//myType[0] = 100
	//
	//fmt.Println(myType)
	//
	//fmt.Println(add(myType))
	//
	//fmt.Println(myType)
	//
	//d := time.Duration(1)
	//
	//var t = time.Now()
	//fmt.Println(t)
	//
	//t = t.Add(d)
	//fmt.Println(t)
}
