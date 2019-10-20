package main

import (
	"time"

	"github.com/go-xorm/xorm"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	_ "github.com/mattn/go-sqlite3"
)

/*
   go get -u github.com/mattn/go-sqlite3
   go get -u github.com/go-xorm/xorm
   如果您使用的是win64并且无法安装go-sqlite3：
       1.下载：https：//sourceforge.net/projects/mingw-w64/files/latest/download
       2.选择“x86_x64”和“posix”
       3.添加C:\Program Files\mingw-w64\x86_64-7.1.0-posix-seh-rt_v5-rev1\mingw64\bin
       到你的PATH env变量。
   手册: http://xorm.io/docs/
*/

//User是我们的用户表结构。
type User struct {
	ID        int64     // xorm默认自动递增
	Username  string    `xorm:"varchar(200)"`
	Email     string    `xorm:"varchar(100)"`
	Password  string    `xorm:"varchar(200)"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

func main() {
	// 数据库时区问题
	time.LoadLocation("Asia/Shanghai")

	app := iris.New()
	// 默认日志
	// app.Logger().SetLevel("debug")

	// app.Use(recover.New())
	// app.Use(logger.New())

	// 自定义日志
	customLogger := logger.New(logger.Config{
		//状态显示状态代码
		Status: true,
		// IP显示请求的远程地址
		IP: true,
		//方法显示http方法
		Method: true,
		// Path显示请求路径
		Path: true,
		// Query将url查询附加到Path。
		Query: true,
		//Columns：true，
		// 如果不为空然后它的内容来自`ctx.Values(),Get("logger_message")
		//将添加到日志中。
		MessageContextKeys: []string{"logger_message"},
		//如果不为空然后它的内容来自`ctx.GetHeader（“User-Agent”）
		MessageHeaderKeys: []string{"User-Agent"},
	})
	app.Use(customLogger)

	//因此，http错误有自己的处理程序
	//注册中间人应该手动完成。
	/*
	   app.OnErrorCode(404 ,customLogger, func(ctx iris.Context) {
	      ctx.Writef("My Custom 404 error page ")
	   })
	*/
	//或捕获所有http错误:
	app.OnAnyErrorCode(customLogger, func(ctx iris.Context) {
		//这应该被添加到日志中，因为`logger.Config＃MessageContextKey`
		ctx.Values().Set("logger_message", "a dynamic message passed to the logs")
		ctx.Writef("My Custom error page")
	})

	orm, err := xorm.NewEngine("sqlite3", "test.db")
	// orm, err := xorm.NewEngine("mysql", "username:password@tcp(host:3306)/dbname?charset=utf8")
	// orm, err := xorm.NewEngine("postgres", "user=test password=test123 dbname=testdb host=127.0.0.1 port=5432 sslmode=disable")
	if err != nil {
		app.Logger().Fatalf("orm failed to initialized: %v", err)
	}
	iris.RegisterOnInterrupt(func() {
		orm.Close()
	})
	err = orm.Sync2(new(User))
	if err != nil {
		app.Logger().Fatalf("orm failed to initialized User table: %v", err)
	}
	app.Post("/create", func(ctx iris.Context) {
		user := &User{}
		if err := ctx.ReadJSON(user); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(err.Error())
			return
		}
		orm.Insert(user)
		ctx.Writef("user inserted: %#v", user)
	})
	app.Get("/get/{id:int}", func(ctx iris.Context) {
		id, _ := ctx.Params().GetInt("id")
		//int到int64
		id64 := int64(id)
		ctx.Writef("id is %#v", id64)

		user := User{ID: id64}
		if ok, _ := orm.Get(&user); ok {
			ctx.Writef("user found: %#v", user)
		}
	})
	app.Delete("/delete/{id:int}", func(ctx iris.Context) {
		id, _ := ctx.Params().GetInt("id")
		//int到int64
		id64 := int64(id)
		user := User{ID: id64}
		orm.Delete(user)
		ctx.Writef("user delete: %#v", user)
	})
	app.Put("/update/{id:int}", func(ctx iris.Context) {
		id, _ := ctx.Params().GetInt("id")
		//int到int64
		id64 := int64(id)

		user := &User{}
		if err := ctx.ReadJSON(user); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(err.Error())
			return
		}

		orm.Id(id64).Update(user)
		ctx.Writef("user update: %#v", user)
	})
	// http://localhost:8080/create
	// http://localhost:8080/get/id:int
	// http://localhost:8080/delete/id:int
	// http://localhost:8080/update/id:int
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
