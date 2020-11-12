package api

import (
	"github.com/kataras/iris/v12"
)

type User struct{}

func (u *User) GetUserinfo(ctx iris.Context) {
	ctx.Text("我是一个用户登录操作handler")
}
