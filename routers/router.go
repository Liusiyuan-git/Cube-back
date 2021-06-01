package routers

import (
	"Cube-back/controllers/community"
	"Cube-back/controllers/login"
	"Cube-back/controllers/register"
	"Cube-back/controllers/user"
	"Cube-back/log"
	beego "github.com/beego/beego/v2/server/web"
)

func apiRegister() {
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/register",
			beego.NSRouter("/user.register", &register.Controller{}, "post:UserRegister"),
		),
		beego.NSNamespace("/user",
			beego.NSRouter("/verification.code", &user.Controller{}, "post:VerificationCode"),
			beego.NSRouter("/password.change", &user.Controller{}, "post:PasswordChange"),
			beego.NSRouter("/get.photo", &user.Controller{}, "get:GetUserPhoto"),
		),
		beego.NSNamespace("/login",
			beego.NSRouter("/user.login", &login.Controller{}, "post:UserLogin"),
			beego.NSRouter("/count.exit", &login.Controller{}, "post:CountExit"),
			beego.NSRouter("/login.status", &login.Controller{}, "post:LoginStatus"),
		),
		beego.NSNamespace("/community",
			beego.NSRouter("/article.like", &community.Controller{}, "post:Test"),
		),
	)
	beego.AddNamespace(ns)
}

func init() {
	apiRegister()
	log.Info("router init successfully")
}
