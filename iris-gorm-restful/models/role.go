package models

import (
	"../database"
	"github.com/jinzhu/gorm"
	"github.com/kataras/golog"
)

type Role struct {
	gorm.Model
	Name        string `gorm:"unique;not null VARCHAR(191)"`
	DisplayName string `gorm:"VARCHAR(191)"`
	Desc        string `gorm:"VARCHAR(191)"`
}

type RoleJson struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Desc        string `json:"desc"`
}

/**
 * 通过 id 获取 role 记录
 * @method GetRoleById
 * @param  {[type]}       role  *Role [description]
 */
func GetRoleById(id uint) (role *Role, err error) {
	role = new(Role)
	role.ID = id

	if err = database.DB.First(role).Error; err != nil {
		golog.Error("GetRoleByIdErr ", err)
	}

	return
}

/**
 * 通过 name 获取 role 记录
 * @method GetRoleByName
 * @param  {[type]}       role  *Role [description]
 */
func GetRoleByName(name string) (role *Role, err error) {
	role = new(Role)
	role.Name = name

	if err = database.DB.First(&role).Error; err != nil {
		golog.Error("GetRoleByNameErr ", err)
	}

	return
}

/**
 * 通过 id 删除角色
 * @method DeleteRoleById
 */
func DeleteRoleById(id uint) {
	u := new(Role)
	u.ID = id

	if err := database.DB.Delete(u).Error; err != nil {
		golog.Error("DeleteRoleErr ", err)
	}
}

/**
 * 获取所有的角色
 * @method GetAllRole
 * @param  {[type]} name string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllRoles(name, orderBy string, offset, limit int) (roles []*Role, err error) {

	if err = database.GetAll(name, orderBy, offset, limit).Find(&roles).Error; err != nil {
		golog.Error("GetAllRoleErr ", err)
	}
	return
}

/**
 * 创建
 * @method CreateRole
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func CreateRole(aul *RoleJson) (role *Role, err error) {

	role = new(Role)
	role.Name = aul.Name
	role.DisplayName = aul.DisplayName
	role.Desc = aul.Desc

	if err = database.DB.Create(role).Error; err != nil {
		golog.Error("CreateRoleErr ", err)
	}

	// perms := []Permission{}
	// database.DB.Where("id in (?)", permIds).Find(&perms)
	// golog.Error(perms)
	// if err := database.DB.Model(&role).Association("Perms").Append(perms).Error; err != nil {
	// 	golog.Error("AppendPermsErr ", err)
	// }

	return
}

/**
 * 更新
 * @method UpdateRole
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func UpdateRole(rj *RoleJson, id uint) (role *Role, err error) {
	role = new(Role)
	role.ID = id

	if err = database.DB.Model(&role).Updates(rj).Error; err != nil {
		golog.Error("UpdatRoleErr ", err)
	}

	// perms := []Permission{}
	// database.DB.Where("id in (?)", permIds).Find(&perms)
	// if err := database.DB.Model(&role).Association("Perms").Replace(perms).Error; err != nil {
	// 	golog.Error("AppendPermsErr ", err)
	// }

	return
}

/**
*创建系统管理员
*@return   *models.AdminRoleTranform api格式化后的数据格式
 */
func CreateSystemAdminRole(rolename string) (role *Role, err error) {
	aul := new(RoleJson)
	aul.Name = rolename
	aul.DisplayName = "超级管理员"
	aul.Desc = "超级管理员"

	role, err = GetRoleByName(aul.Name)

	if role.ID == 0 {
		golog.Info("创建角色")
		return CreateRole(aul)
	} else {
		golog.Warn("角色已存在")
		return
	}
}
