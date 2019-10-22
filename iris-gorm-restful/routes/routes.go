package routes

import (
	"../controllers"
	"../middleware"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/router"
)

func Register(api *iris.Application) {
	// allow cors
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowedMethods:   []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
	})

	main := api.Party("/", crs).AllowMethods(iris.MethodOptions)

	home := main.Party("/")
	home.Get("/", func(ctx iris.Context) { // 首页模块
		_ = ctx.View("index.html")
	})
	mainDoc := api.Party("/apidoc", crs).AllowMethods(iris.MethodOptions)
	mainDoc.Get("/", func(ctx iris.Context) { // 首页模块
		_ = ctx.View("/apidoc/index.html")
	})

	v1 := api.Party("/v1", crs).AllowMethods(iris.MethodOptions)
	{
		v1.Use(middleware.NewYaag()) // <- IMPORTANT, register the middleware.
		v1.Post("/api/login", controllers.UserLogin)
		v1.PartyFunc("/api", func(admin router.Party) {
			admin.Use(middleware.JwtHandler().Serve, middleware.AuthToken)
			admin.Get("/logout", controllers.UserLogout)

			admin.PartyFunc("/users", func(users router.Party) {
				users.Get("/", controllers.GetAllUsers)
				users.Get("/{id:uint}", controllers.GetUser)
				users.Post("/", controllers.CreateUser)
				users.Put("/{id:uint}", controllers.UpdateUser)
				users.Delete("/{id:uint}", controllers.DeleteUser)
				users.Get("/profile", controllers.GetProfile)
				users.Patch("/changePasswd/{id:uint}", controllers.UpdateUserPassword)
			})
			admin.PartyFunc("/roles", func(roles router.Party) {
				roles.Get("/", controllers.GetAllRoles)
				roles.Get("/{id:uint}", controllers.GetRole)
				roles.Post("/", controllers.CreateRole)
				roles.Put("/{id:uint}", controllers.UpdateRole)
				roles.Delete("/{id:uint}", controllers.DeleteRole)
			})
		})
	}
}