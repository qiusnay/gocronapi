package main

import (
	"github.com/kataras/iris/v12"
	"github.com/qiusnay/gocronapi/route"
	"github.com/qiusnay/gocronapi/util"
)

func main() {
	app := iris.New()
	util.SetBASEPath()
	//路由注册
	route.RegisterRouter(app)
	// 设置静态资源
	// app.StaticWeb("/static", "static")
	// app.Run(iris.Addr(":8088"), iris.WithConfiguration(util.GetMyConf()))
	conf := util.LoadConf()
	app.Run(iris.Addr(":8080"), iris.WithConfiguration(conf))
}
