// 用户
package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"

	"step/models"
	"fmt"
)

type UserController struct {
	beego.Controller
}

// Index 用户列表
func (u *UserController) Index() {
	u.Data["users"] = models.List()
	u.Data["status"] = models.Status()
	u.Data["roles"] = models.Roles()
	u.Layout = "base.html"
	u.TplName = "user/index.html"
}

// Create 添加用户
func (u *UserController) Create() {
	beego.ReadFromRequest(&u.Controller)
	if u.Ctx.Request.Method == "POST" {
		nickname := u.Input().Get("nickname")
		email := u.Input().Get("email")
		password := u.Input().Get("password")
		role, err := strconv.Atoi(u.Input().Get("role"))
		if err != nil {
			u.Redirect("/user/create", 302)
		}
		// 用户名成不能为空
		if strings.TrimSpace(nickname) == "" {
			u.Redirect("/user/create", 302)
		}
		// 邮箱不能为空
		if strings.TrimSpace(email) == "" {
			u.Redirect("/user/create", 302)
		}
		// 邮箱格式不合法
		// 密码不能为空
		if strings.TrimSpace(password) == "" {
			u.Redirect("/user/create", 302)
		}
		var id int64
		if id, err = models.Save(&models.User{
			Nickname: nickname,
			Email:    email,
			Password: password,
			Role:     role,
		}); err != nil {
			u.Redirect("/user/create", 302)
		}
		u.Redirect("/user/" + strconv.FormatInt(id, 10), 302)
	}
	u.Data["roles"] = models.Roles()
	u.Data["user"] = models.User{}
	u.Data["isNewRecord"] = true
	u.Layout = "base.html"
	u.TplName = "user/create.html"
}

// View 用户查看
func (u *UserController) View() {
	id, err := strconv.Atoi(u.Ctx.Input.Param(":id"))
	if err != nil || id <= 0 {
		u.Redirect("/user", 302)
	}
	user, err := models.Info(int64(id))
	if err != nil {
		u.Redirect("/user", 302)
	}
	u.Data["user"] = user
	u.Data["role"] = models.RoleDesc(user.Role)
	u.Data["status"] = models.StatusDesc(user.Status)
	u.Data["ctime"] = time.Unix(user.Ctime, 0).Format("2006-01-02 15:04:05")
	u.Data["utime"] = time.Unix(user.Utime, 0).Format("2006-01-02 15:04:05")
	u.Layout = "base.html"
	u.TplName = "user/view.html"
}

// Update 用户编辑
func (u *UserController) Update() {
	id, err := strconv.Atoi(u.Ctx.Input.Param(":id"))
	if err != nil || id <= 0 {
		u.Redirect("/user", 302)
	}
	user, err := models.Info(int64(id))
	if err != nil {
		u.Redirect("/user", 302)
	}
	if u.Ctx.Request.Method == "POST" {
		nickname := u.Input().Get("nickname")
		email := u.Input().Get("email")
		password := u.Input().Get("password")
		role, err := strconv.Atoi(u.Input().Get("role"))
		if err != nil {
			u.Redirect("/user/update/"+strconv.FormatInt(int64(id), 10), 302)
		}
		user.Nickname = strings.TrimSpace(nickname)
		user.Email = strings.TrimSpace(email)
		user.Password = strings.TrimSpace(password)
		user.Role = role
		var newId int64
		if newId, err = models.Update(&user); err != nil {
			u.Redirect("/user/update/"+strconv.FormatInt(newId, 10), 302)
		}
		u.Redirect("/user/"+strconv.FormatInt(newId, 10), 302)
	}

	u.Data["user"] = user
	u.Data["isNewRecord"] = false
	u.Data["roles"] = models.Roles()
	u.Layout = "base.html"
	u.TplName = "user/update.html"
}

// Delete 删除用户
func (u *UserController) Delete() {
	id, err := strconv.Atoi(u.Ctx.Input.Param(":id"))
	if err != nil || id <= 0 {
		u.Redirect("/user", 302)
	}
	user, err := models.Info(int64(id))
	if err != nil {
		u.Redirect("/user", 302)
	}
	user.Status = 2
	if _, err := models.Update(&user); err != nil {
		fmt.Println(err)
	}
	u.Redirect("/user", 302)
}
