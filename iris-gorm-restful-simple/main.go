package main

import (
	"os"

	"./database"
	"./models"

	"flag"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func main() {
	flag.Parse()
	app := newApp()

	err := app.Run(iris.Addr(":8000"), iris.WithoutServerError(iris.ErrServerClosed))
	if err != nil {
		panic(err)
	}
}

func newApp() *iris.Application {
	hostname, _ := os.Hostname()
	var env string
	if hostname == "wahaha" {
		env = "prod"
		golog.Info("进入线上环境")
	} else {
		env = "test"
		golog.Info("进入测试环境")
	}

	app := iris.New()
	app.Configure(iris.WithOptimizations)

	app.Use(recover.New())
	app.Use(logger.New())

	// app.Logger().SetLevel("debug")
	// app.Use(logger.New())
	// app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
	// 	ctx.JSON(controllers.ApiResource(false, nil, "404 Not Found"))
	// })
	app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		ctx.WriteString("Oups something went wrong, try again")

	})

	// migrate db
	database.DB.AutoMigrate(
		&models.User{},
		&models.Role{},
	)

	iris.RegisterOnInterrupt(func() {
		_ = database.DB.Close()
	})

	// allow cors
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	})
	app.Use(crs)
	app.AllowMethods(iris.MethodOptions)

	//初始化系统 账号 权限 角色
	models.CreateSystemData(env)

	return app
}
