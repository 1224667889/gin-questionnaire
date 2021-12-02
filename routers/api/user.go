package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/models"
	"main/pkg/app"
	"main/pkg/e"
	"main/pkg/util"
	"net/http"
)

// LoginAPI 用户登录
// 输出账号密码进行登录 账号也可使用邮箱
func LoginAPI(c *gin.Context) {
	g := app.Gin{C: c}
	params := struct {
		Account  string `form:"account" json:"account" xml:"account" binding:"required"`
		Password string `form:"password" json:"password" xml:"password" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&params); err != nil {
		g.Response(http.StatusOK, e.INVALID_PARAMS, err.Error())
		return
	}
	var user models.Account
	if err := models.First(&user,
		fmt.Sprintf("WHERE account = '%s' OR email = '%s'", params.Account, params.Account)); err != nil {
		g.Response(http.StatusOK, e.INVALID_PARAMS, "账号不存在")
		return
	}
	if user.Password != params.Password {
		g.Response(http.StatusOK, e.INVALID_PARAMS, "密码错误")
		return
	}
	if token, err := util.GenerateToken(user.Id); err != nil {
		g.Response(http.StatusOK, e.ERROR, err.Error())
	} else {
		g.Response(http.StatusOK, e.SUCCESS, token)
	}
	return
}

// UpdateUser 更新用户信息
// 修改用户个人信息
func UpdateUser(c *gin.Context) {
	g := app.Gin{C: c}
	params := struct {
		Name  string `form:"name" json:"name" xml:"name"`
		Email string `form:"email" json:"email" xml:"email"`
	}{}
	if err := c.ShouldBindJSON(&params); err != nil {
		g.Response(http.StatusOK, e.INVALID_PARAMS, "参数错误")
		return
	}
	id := c.MustGet("id").(string)
	user := models.Account{Id: id}
	if err := models.FindByKey(&user); err != nil {
		g.Response(http.StatusOK, e.ERROR_AUTH_CHECK_TOKEN_FAIL, "用户信息未找到")
		return
	}
	user.Name = params.Name
	user.Email = params.Email
	if err := models.Update(&user); err != nil {
		g.Response(http.StatusOK, e.ERROR_DB, err.Error())
		return
	}
	g.Response(http.StatusOK, e.SUCCESS, "信息修改成功")
	return
}

// ResetPassword 修改密码
// 通过验证原密码修改新密码
func ResetPassword(c *gin.Context) {
	g := app.Gin{C: c}
	params := struct {
		OldPassword string `form:"old_password" json:"old_password" xml:"old_password" binding:"required"`
		NewPassword string `form:"new_password" json:"new_password" xml:"new_password" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&params); err != nil {
		g.Response(http.StatusOK, e.INVALID_PARAMS, "参数错误")
		return
	}
	if params.OldPassword == params.NewPassword {
		g.Response(http.StatusOK, e.INVALID_PARAMS, "新旧密码不能一致")
		return
	}
	id := c.MustGet("id").(string)
	user := models.Account{Id: id}
	if err := models.FindByKey(&user); err != nil {
		g.Response(http.StatusOK, e.ERROR_AUTH_CHECK_TOKEN_FAIL, "用户信息未找到")
		return
	}
	if user.Password != params.OldPassword {
		g.Response(http.StatusOK, e.INVALID_PARAMS, "密码验证失败")
		return
	}
	user.Password = params.NewPassword
	if err := models.Update(&user); err != nil {
		g.Response(http.StatusOK, e.ERROR_DB, err.Error())
		return
	}
	g.Response(http.StatusOK, e.SUCCESS, "密码修改成功")
	return
}

// RegisterAPI 用户注册
// 注册新用户 权限为User
func RegisterAPI(c *gin.Context) {
	g := app.Gin{C: c}
	params := struct {
		Account  string `form:"account" json:"account" xml:"account" binding:"required"`
		Password string `form:"password" json:"password" xml:"password" binding:"required"`
		Name     string `form:"name" json:"name" xml:"name" binding:"required"`
		Email    string `form:"email" json:"email" xml:"email" binding:"email"`
	}{}
	if err := c.ShouldBindJSON(&params); err != nil {
		g.Response(http.StatusOK, e.INVALID_PARAMS, "参数错误")
		return
	}
	user := models.Account{
		Name:     params.Name,
		Account:  params.Account,
		Email:    params.Email,
		Password: params.Password,
		RoleId:   "3",
	}
	if err := models.Insert(&user); err != nil {
		g.Response(http.StatusOK, e.ERROR_DB, "该信息已被其他用户绑定")
		return
	}
	g.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"account": user.Account,
		"email": user.Email,
		"name": user.Name,
	})
	return
}