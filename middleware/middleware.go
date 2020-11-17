package middleware

import (
	"strings"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/qiusnay/gocronapi/middleware/jwt"
)

func CorsHandler() iris.Handler {
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, //允许通过的主机名称
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		Debug:          false,
	})
	return crs
}

func JwtAuthHandler(ctx iris.Context) {
	path := ctx.Path()
	//登录与静态资源页直接跳过
	if strings.Contains(path, "/static") || strings.Contains(path, "/login") {
		ctx.Next()
		return
	}
	// jwt token拦截
	if !jwt.Serve(ctx) {
		return
	}
	ctx.Next()
}
