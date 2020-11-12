package main

import (
	"github.com/kataras/iris/v12"
	"github.com/qiusnay/gocronadmin/route"
)

func main() {
	app := iris.New()
	//路由注册
	route.RegisterRouter(app)
	// 设置静态资源
	// app.StaticWeb("/static", "static")
	// app.Run(iris.Addr(":8088"), iris.WithConfiguration(util.GetMyConf()))
	app.Run(iris.Addr(":8080"), iris.WithConfiguration(iris.YAML("./conf/app.yml")))
}
