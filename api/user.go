package api

import "fmt"

type User struct{}

func (u *User) Login() {
	fmt.Println("我是一个用户登录操作handler")
}
