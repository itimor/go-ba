package main

import (
	"./config"
	"./database"
	"./models"
	"./routes"

	"flag"
	"os"

	"github.com/betacraft/yaag/yaag"
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
	loglevle := config.Conf.Get("test.loglevel").(string)
	app.Logger().SetLevel(loglevle)
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
	golog.Info("初始化数据库")
	database.DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.OauthToken{},
	)

	iris.RegisterOnInterrupt(func() {
		_ = database.DB.Close()
	})

	// 加载路由
	routes.Register(app)

	//api 文档配置
	appName := config.Conf.Get("server.name").(string)
	appDoc := config.Conf.Get("server.apidoc").(string)
	appURL := config.Conf.Get("server.apiurl").(string)
	yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware.
		On:       true,
		DocTitle: appName,
		DocPath:  appDoc + "/index.html", //设置绝对路径
		BaseUrls: map[string]string{
			"Production": appURL,
			"Staging":    "",
		},
	})

	//初始化系统 账号 权限 角色
	models.CreateSystemData(env)

	return app
}
