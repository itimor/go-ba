package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

type Company struct {
	Id    int64
	Name  string
	City  string
	Other string
}

func main() {
	app := iris.New()

	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())

	app.Post("/", MyHandler)
	app.Post("/slice", MyHandler2)
	//使用Postman或其他什么来做POST请求
	//使用RAW BODY到http//localhost：8080：
	/*
	   {
	       "Name": "iris-Go",
	       "City": "New York",
	       "Other": "Something here"
	   }
	*/
	//和Content-Type到application/json（可选且易用）
	//响应应该是：
	//接收值: main.Company{Name:"iris-Go", City:"New York", Other:"Something here"}
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
}

func MyHandler(ctx iris.Context) {
	// var c Company   // use &c
	c := &Company{}
	if err := ctx.ReadJSON(c); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
	ctx.Writef("Received: %#+v\n", c)
}

// 简单的json，请阅读https://golang.org/pkg/encoding/json
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// MyHandler2从JSON POST数据中读取Person的集合。
func MyHandler2(ctx iris.Context) {
	var persons []Person
	err := ctx.ReadJSON(&persons)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}
	ctx.Writef("Received: %#+v\n", persons)
}
