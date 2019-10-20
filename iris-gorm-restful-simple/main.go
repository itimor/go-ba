package main

import (
	"./conf"

	"flag"
	"fmt"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
)

func main() {
	flag.Parse()
	app := newApp()

	// 读取“./config.json”配置文件
	fmt.Println("Port", conf.Sysconfig.App.Port)
	fmt.Println("JWTTimeout", conf.Sysconfig.App.JWTTimeout)
	fmt.Println("Charset", conf.Sysconfig.DB.Charset)

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}

func newApp() *iris.Application {
	app := iris.New()
	app.Configure(iris.WithOptimizations)
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})
	app.Use(crs)
	app.AllowMethods(iris.MethodOptions)
	return app
}
