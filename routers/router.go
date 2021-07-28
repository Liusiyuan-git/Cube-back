package routers

import (
	"Cube-back/controllers/blog"
	"Cube-back/controllers/common"
	"Cube-back/controllers/draft"
	"Cube-back/controllers/login"
	"Cube-back/controllers/register"
	"Cube-back/controllers/talk"
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
		beego.NSNamespace("/private",
			beego.NSRouter("/send.blog", &blog.Controller{}, "post:BlogSend"),
			beego.NSRouter("/send.draft", &draft.Controller{}, "post:DraftSend"),
			beego.NSRouter("/get.draft", &draft.Controller{}, "post:DraftGet"),
			beego.NSRouter("/remove.draft", &draft.Controller{}, "post:DraftRemove"),
			beego.NSRouter("/blog.collect", &blog.Controller{}, "post:BlogCollect"),
			beego.NSRouter("/blog.comment.send", &blog.Controller{}, "post:BlogCommentSend"),
			beego.NSRouter("/blog.collect.confirm", &blog.Controller{}, "post:BlogCollectConfirm"),
			beego.NSRouter("/send.talk", &talk.Controller{}, "post:TalkSend"),
			beego.NSRouter("/send.talk.comment", &talk.Controller{}, "post:TalkCommentSend"),
			beego.NSRouter("/delete.talk.Comment", &talk.Controller{}, "post:TalkCommentDelete"),
			beego.NSRouter("/cube.collection.get", &blog.Controller{}, "post:BlogCollectionGet"),
		),
		beego.NSNamespace("/common",
			beego.NSRouter("/blog.get", &common.Controller{}, "post:BlogGet"),
			beego.NSRouter("/blog.detail", &common.Controller{}, "post:BlogDetail"),
			beego.NSRouter("/blog.like", &common.Controller{}, "post:BlogLike"),
			beego.NSRouter("/blog.comment.get", &common.Controller{}, "post:BlogCommonGet"),
			beego.NSRouter("/talk.get", &common.Controller{}, "post:TalkGet"),
			beego.NSRouter("/talk.comment.get", &common.Controller{}, "post:TalkCommentGet"),
			beego.NSRouter("/talk.like", &common.Controller{}, "post:TalkLike"),
			beego.NSRouter("/cube.information.get", &common.Controller{}, "post:CubeInformationGet"),
			beego.NSRouter("/cube.view.get", &common.Controller{}, "post:CubeViewGet"),
		),
	)
	beego.AddNamespace(ns)
}

func init() {
	apiRegister()
	log.Info("router init successfully")
}
