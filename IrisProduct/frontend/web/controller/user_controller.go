package controller

import (
	"../../../services"
	"../../../datamodels"
	"../../../tool"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
	"strconv"
	"../../../encrypt"
)

type UserController struct {
	Ctx         iris.Context
	UserService services.IUserService
	Session     *sessions.Session
}

func (c *UserController) GetRegister() mvc.View {
	return mvc.View{
		Name: "user/register.html",
	}
}

func (c *UserController) PostRegister() {
	var (
		nickName = c.Ctx.FormValue("nickName")
		userName = c.Ctx.FormValue("userName")
		pwd      = c.Ctx.FormValue("password")
	)
	user := &datamodels.User{
		NickName: nickName,
		UserName: userName,
		Password: pwd,
	}
	_, err := c.UserService.AddUser(user)
	fmt.Println(err)
	if err != nil {
		c.Ctx.Redirect("user/error")
		return
	}
	c.Ctx.Redirect("login")
	return
}

func (c *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name: "user/login.html",
	}
}

func (c *UserController) PostLogin() mvc.Response {
	user := &datamodels.User{}
	c.Ctx.ReadForm(user)
	dbUser, isOk := c.UserService.IsLoginSuccess(user.UserName, user.Password)
	if !isOk {
		return mvc.Response{
			Path: "login",
		}
	}
	fmt.Println(user, dbUser)
	uid := strconv.FormatInt(dbUser.ID, 10)
	tool.GlobalCookie(c.Ctx, "uid", uid)

	enuid, err := encrypt.EnPwdCode(uid)
	if err != nil {
		c.Ctx.Application().Logger().Debug(err)
	}
	//2.写入用户ID到cookie中
	tool.GlobalCookie(c.Ctx, "sign", enuid)
	//c.Session.Set("userID", strconv.FormatInt(dbUser.ID, 10))
	return mvc.Response{
		Path: "/product/detail",
	}
}
