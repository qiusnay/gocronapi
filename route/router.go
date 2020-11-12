package route

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/qiusnay/gocronadmin/api"
)

func RegisterRouter(app *iris.Application) {
	/* 定义路由 */
	main := app.Party("/", CorsHandler()).AllowMethods(iris.MethodOptions)
	main.Use(JwtAuthHandler)

	home := main.Party("/")
	home.Get("/", func(ctx iris.Context) { // 首页模块
		ctx.HTML("<h1>Welcome</h1>")
	})

	mvc.Configure(app.Party("/user"), func(app *mvc.Application) {
		app.Handle(new(api.User))
	})

	// 用户API模块
	// user := main.Party("/user")
	// app.Handle("GET", "/user", new(api.User))
	// user.Get("/userinfo", api.User{})

	// 权限API模块
	// admin := main.Party("/admin")
	// {
	// 	// 用户管理
	// 	users := admin.Party("/users")
	// 	users.Get("/", dispatch.Handler(routes.UserTable))                   // 用户列表
	// 	users.Put("/", dispatch.Handler(routes.UpdateUser))                  // 更新用户
	// 	users.Delete("/{uids:string}", dispatch.Handler(routes.DeleteUsers)) // 删除用户

	// 	//角色管理
	// 	role := admin.Party("/role")
	// 	role.Get("/", dispatch.Handler(routes.RoleTable))                       // 角色报表
	// 	role.Put("/", dispatch.Handler(routes.UpdateRole))                      // 更新角色
	// 	role.Post("/", dispatch.Handler(routes.CreateRole))                     // 创建角色
	// 	role.Delete("/{rids:string}", dispatch.Handler(routes.DeleteRole))      // 删除角色
	// 	role.Get("/user/{rKey:string}", dispatch.Handler(routes.RoleUserTable)) // 角色关联的用户表
	// 	role.Get("/menu/{rid:long}", dispatch.Handler(routes.RoleMenuTable))    // 角色关联的菜单表

	// 	//菜单管理
	// 	menu := admin.Party("/menu")
	// 	menu.Get("/", dispatch.Handler(routes.MenuTable))                    // 菜单列表
	// 	menu.Post("/permissions", dispatch.Handler(routes.RelationUserRole)) // 给角色添加权限
	// }

	// demo测试API模块
	// demo := main.Party("/demo")
	// {
	// 	demo.Get("/{pid:long}", dispatch.Handler(routes.GetOneProduct))
	// 	demo.Put("/", dispatch.Handler(routes.AddOneProduct))
	// }
}
