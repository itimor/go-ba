package main

import (
	"fmt"

	"./config"

	"flag"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
)

func main() {
	flag.Parse()
	app := newApp()

	server_port := config.Tomlconfig.Get("app.server_port").(string)
	fmt.Printf("format:", server_port)
	app.Run(iris.Addr(":"+server_port), iris.WithoutServerError(iris.ErrServerClosed))
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
