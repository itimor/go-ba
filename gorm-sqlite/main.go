package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

var db *gorm.DB

// User 是我们的用户表结构。
type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(100);not null;unique"`
	Email    string `gorm:"type:varchar(20);not null;unique"`
	Password string `gorm:"not null"`
}

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "test.db")
	// db, err = gorm.Open("postgres", "user=gorm password=gorm DB.name=gorm port=9920 sslmode=disable")
	// db, err = gorm.Open("mysql", "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True")
	// db, err = gorm.Open("mssql", "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm")
	if err != nil {
		panic("Failed to Connect to Database")
	}
	db.LogMode(true)
}

func main() {
	// 初始化连接数据库
	// db.DropTable(&User{})
	db.AutoMigrate(&User{})

	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())

	app.Post("/create", func(ctx iris.Context) {
		user := &User{}
		if err := ctx.ReadJSON(user); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(err.Error())
			return
		}
		db.Create(user)
		ctx.Writef("user inserted: %v", user)
	})
	app.Get("/get/{id:int}", func(ctx iris.Context) {
		id, _ := ctx.Params().GetInt("id")
		//int到int64
		id64 := int64(id)

		user := db.First(&User{}, id64)
		ctx.Writef("user found: %v", user)
	})
	app.Delete("/delete/{id:int}", func(ctx iris.Context) {
		id, _ := ctx.Params().GetInt("id")
		//int到int64
		id64 := int64(id)
		// 软删除，记录DeletedAt,DeletedAt默认为null
		// user := db.Delete(&User{}, "id=?", id64)
		// 硬删除
		user := db.Unscoped().Delete(&User{}, "id=?", id64)
		ctx.Writef("user delete: %v", user)
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
		db.Model(&User{}).Where("id=?", id64).Update(user)
		ctx.Writef("user update: %v", user)
	})
	// http://localhost:8080/create
	// http://localhost:8080/get/id:int
	// http://localhost:8080/delete/id:int
	// http://localhost:8080/update/id:int
	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
}
