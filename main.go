package main

import (
	"github.com/kataras/iris/v12"
	"github.com/qiusnay/gocronadmin/init"
	"github.com/qiusnay/gocronadmin/util"
)

func main() {
	app := iris.New()
	//路由注册
	init.RegisterRouter(app)
	// 设置静态资源
	app.StaticWeb("/static", "static")
	app.Run(iris.Addr(":8088"), iris.WithConfiguration(util.GetMyConf()))
}
