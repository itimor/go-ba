package controllers

import (
	"fmt"

	"../models"

	"github.com/kataras/iris"
	"gopkg.in/go-playground/validator.v9"
)

func GetProfile(ctx iris.Context) {
	uid := ctx.Values().Get("uid").(uint)
	user, _ := models.GetUserById(uid)

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, user, "success"))
}

func GetUser(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	user, _ := models.GetUserById(id)

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, user, "success"))
}

func CreateUser(ctx iris.Context) {

	aul := new(models.UserJson)

	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		_, _ = ctx.JSON(errorData(err))
	} else {
		err := validate.Struct(aul)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			for _, err := range err.(validator.ValidationErrors) {
				fmt.Println()
				fmt.Println(err.Namespace())
				fmt.Println(err.Field())
				fmt.Println(err.Type())
				fmt.Println(err.Param())
				fmt.Println()
			}
		} else {
			u, _ := models.CreateUser(aul)

			ctx.StatusCode(iris.StatusOK)
			if u.ID == 0 {
				_, _ = ctx.JSON(ApiResource(false, u, "error"))
			} else {
				_, _ = ctx.JSON(ApiResource(true, u, "success"))
			}
		}
	}
}

func UpdateUser(ctx iris.Context) {
	aul := new(models.UserJson)

	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		_, _ = ctx.JSON(errorData(err))
	} else {
		err := validate.Struct(aul)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			for _, err := range err.(validator.ValidationErrors) {
				fmt.Println()
				fmt.Println(err.Namespace())
				fmt.Println(err.Field())
				fmt.Println(err.Type())
				fmt.Println(err.Param())
				fmt.Println()
			}
		} else {
			id, _ := ctx.Params().GetInt("id")
			uid := uint(id)

			u, _ := models.UpdateUser(aul, uid)
			ctx.StatusCode(iris.StatusOK)
			if u.ID == 0 {
				_, _ = ctx.JSON(ApiResource(false, u, "error"))
			} else {
				_, _ = ctx.JSON(ApiResource(true, u, "success"))
			}
		}
	}
}

func DeleteUser(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	models.DeleteUserById(id)

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, nil, "success"))
}

func GetAllUsers(ctx iris.Context) {
	//offset := utils.Tool.ParseInt(ctx.FormValue("offset"), 1)
	//limit := utils.Tool.ParseInt(ctx.FormValue("limit"), 20)
	//name := ctx.FormValue("name")
	//username := ctx.FormValue("username")
	//orderBy := ctx.FormValue("orderBy")
	offset := ctx.URLParamIntDefault("offset", 1)
	limit := ctx.URLParamIntDefault("limit", 15)
	name := ctx.URLParam("name")
	orderBy := ctx.URLParam("orderBy")

	users := models.GetAllUsers(name, orderBy, offset, limit)

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, users, "success"))
}
