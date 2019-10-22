package models

import (
	"../config"
)

/**
*初始化系统 账号 权限 角色
 */
func CreateSystemData(env string) {
	// perm := CreateSystemAdminPermission() //初始化权限

	// permIds := []uint{perm.ID}
	// role := CreateSystemAdminRole(permIds) //初始化角色

	rolename := config.Conf.Get(env + ".role").(string)
	role, _ := CreateSystemAdminRole(rolename)

	if role.ID != 0 {
		aul := new(UserJson)
		aul.Username = config.Conf.Get(env + ".user").(string)
		aul.Password = config.Conf.Get(env + ".pass").(string)
		CreateSystemAdmin(aul, []string{rolename})
	}
}
