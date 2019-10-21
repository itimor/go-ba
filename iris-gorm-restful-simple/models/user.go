package models

import (
	"../database"

	"fmt"
	"time"

	"github.com/beego/bee/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/gorm"
	"github.com/kataras/golog"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(50);not null;unique"`
	Password string `gorm:"type:varchar(200);not null"`
	Avatar   string
	Roles    []Role `gorm:"many2many:user_roles;"`
	IsActive bool   `gorm:"default:true"`
	IsAdmin  bool   `gorm:"default:false"`
}

type UserJson struct {
	Username string   `json:"username" validate:"required,gte=2,lte=50"`
	Password string   `json:"password" validate:"required,gte=8,lte=200"`
	Avatar   string   `json:"avatar" validate:"required,gte=2,lte=200"`
	Roles    []string `json:"roles" validate:"required"`
}

/**
 * 通过 id 获取 user 记录
 * @method GetUserById
 * @param  {[type]}       user  *User [description]
 */
func GetUserById(id uint) (user *User, err error) {
	user.ID = id

	if err = database.DB.Preload("Role").First(user).Error; err != nil {
		golog.Error("GetUserByIdErr:%s", err)
	}

	return
}

/**
 * 通过 username 获取 user 记录
 * @method GetUserByUserName
 * @param  {[type]}       user  *User [description]
 */
func GetUserByUserName(username string) (user *User, err error) {
	user := &User{Username: username}

	if err := database.DB.Preload("Role").First(user).Error; err != nil {
		golog.Error("GetUserByUserNameErr:%s", err)
	}

	return nil, user
}

/**
 * 通过 id 删除用户
 * @method DeleteUserById
 */
func DeleteUserById(id uint) {
	u := new(User)
	u.ID = id

	if err := database.DB.Delete(u).Error; err != nil {
		golog.Error("DeleteUserByIdErr:%s", err)
	}
}

/**
 * 获取所有的账号
 * @method GetAllUser
 * @param  {[type]} name string [description]
 * @param  {[type]} username string [description]
 * @param  {[type]} orderBy string [description]
 * @param  {[type]} offset int    [description]
 * @param  {[type]} limit int    [description]
 */
func GetAllUsers(name, orderBy string, offset, limit int) (users []*User) {
	if err := database.GetAll(name, orderBy, offset, limit).Preload("Role").Find(&users).Error; err != nil {
		golog.Error("GetAllUserErr:%s", err)
	}
	return
}

/**
 * 创建
 * @method CreateUser
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func CreateUser(aul *UserJson) (user *User) {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(aul.Password, salt)

	user = new(User)
	user.Username = aul.Username
	user.Password = hash
	user.Avatar = aul.Avatar
	user.Roles = aul.Roles

	if err := database.DB.Create(user).Error; err != nil {
		golog.Error("CreateUserErr:%s", err)
	}

	return
}

/**
 * 更新
 * @method UpdateUser
 * @param  {[type]} kw string [description]
 * @param  {[type]} cp int    [description]
 * @param  {[type]} mp int    [description]
 */
func UpdateUser(uj *UserJson, id uint) *User {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(uj.Password, salt)

	user := new(User)
	user.ID = id
	uj.Password = hash
	user.Avatar = aul.Avatar
	user.Roles = aul.Roles

	if err := database.DB.Model(user).Updates(uj).Error; err != nil {
		golog.Error("UpdateUserErr:%s", err)
	}

	return user
}

/**
 * 判断用户是否登录
 * @method CheckLogin
 * @param  {[type]}  id       int    [description]
 * @param  {[type]}  password string [description]
 */
func CheckLogin(username, password string) (response Token, status bool, msg string) {
	user := UserAdminCheckLogin(username)
	if user.ID == 0 {
		msg = "用户不存在"
		return
	} else {
		if ok := bcrypt.Match(password, user.Password); ok {
			token := jwt.New(jwt.SigningMethodHS256)
			claims := make(jwt.MapClaims)
			claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
			claims["iat"] = time.Now().Unix()
			token.Claims = claims
			tokenString, err := token.SignedString([]byte("secret"))

			if err != nil {
				msg = err.Error()
				return
			}

			oauthToken := new(OauthToken)
			oauthToken.Token = tokenString
			oauthToken.UserId = user.ID
			oauthToken.Secret = "secret"
			oauthToken.Revoked = false
			oauthToken.ExpressIn = time.Now().Add(time.Hour * time.Duration(1)).Unix()
			oauthToken.CreatedAt = time.Now()

			response = oauthToken.OauthTokenCreate()
			status = true
			msg = "登陆成功"

			return

		} else {
			msg = "用户名或密码错误"
			return
		}
	}
}

/**
* 用户退出登陆
* @method UserAdminLogout
* @param  {[type]} ids string [description]
 */
func UserAdminLogout(userId uint) bool {
	ot := UpdateOauthTokenByUserId(userId)
	return ot.Revoked
}

/**
*创建系统管理员
*@param role_id uint
*@return   *models.AdminUserTranform api格式化后的数据格式
 */
func CreateSystemAdmin(roleId uint) *User {

	aul := new(UserJson)
	aul.Username = config.Conf.Get("test.LoginUserName").(string)
	aul.Password = config.Conf.Get("test.LoginPwd").(string)
	aul.Name = config.Conf.Get("test.LoginName").(string)
	aul.RoleID = roleId

	user := GetUserByUserName(aul.Username)

	if user.ID == 0 {
		fmt.Println("创建账号")
		return CreateUser(aul)
	} else {
		fmt.Println("重复初始化账号")
		return user
	}
}
