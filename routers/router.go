package routers

import (
	"Cube-back/controllers/blog"
	"Cube-back/controllers/common"
	"Cube-back/controllers/draft"
	"Cube-back/controllers/login"
	"Cube-back/controllers/message"
	"Cube-back/controllers/profile"
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
			beego.NSRouter("draft.image.upload", &draft.Controller{}, "post:DraftImageUpload"),
			beego.NSRouter("draft.image.delete", &draft.Controller{}, "post:DraftImageDelete"),
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
			beego.NSRouter("/send.user.image", &profile.Controller{}, "post:SendUserImage"),
			beego.NSRouter("/user.introduce.send", &profile.Controller{}, "post:UserIntroduceSend"),
			beego.NSRouter("/user.name.send", &profile.Controller{}, "post:UserNameSend"),
			beego.NSRouter("/user.profile.get", &profile.Controller{}, "post:UserProfileGet"),
			beego.NSRouter("/user.image.update", &profile.Controller{}, "post:UserImageUpdate"),
			beego.NSRouter("/user.care.set", &profile.Controller{}, "post:UserCareSet"),
			beego.NSRouter("/user.care.get", &profile.Controller{}, "post:UserCareGet"),
			beego.NSRouter("/user.care.confirm", &profile.Controller{}, "post:UserCareConfirm"),
			beego.NSRouter("/user.care.cancel", &profile.Controller{}, "post:UserCareCancel"),
			beego.NSRouter("/profile.leave.set", &profile.Controller{}, "post:ProfileLeaveSet"),
			beego.NSRouter("/user.message.get", &message.Controller{}, "post:UserMessageGet"),
			beego.NSRouter("/message.profile.get", &message.Controller{}, "post:MessageProfileGet"),
			beego.NSRouter("/user.message.clean", &message.Controller{}, "post:UserMessageClean"),
			beego.NSRouter("/user.profile.care", &common.Controller{}, "post:UserProfileCare"),
			beego.NSRouter("/message.profile.user.talk.get", &message.Controller{}, "post:MessageProfileUserTalkGet"),
			beego.NSRouter("/message.profile.user.talk.clean", &message.Controller{}, "post:MessageProfileUserTalkClean"),
			beego.NSRouter("/message.profile.user.blog.get", &message.Controller{}, "post:MessageProfileUserBlogGet"),
			beego.NSRouter("/message.profile.user.blog.clean", &message.Controller{}, "post:MessageProfileUserBlogClean"),
			beego.NSRouter("/blog.delete", &blog.Controller{}, "post:BlogDelete"),
			beego.NSRouter("/collect.delete", &blog.Controller{}, "post:CollectDelete"),
		),
		beego.NSNamespace("/common",
			beego.NSRouter("/blog.get", &common.Controller{}, "post:BlogGet"),
			beego.NSRouter("/forum.blog.get", &common.Controller{}, "post:BlogForumGet"),
			beego.NSRouter("/blog.detail", &common.Controller{}, "post:BlogDetail"),
			beego.NSRouter("/blog.like", &common.Controller{}, "post:BlogLike"),
			beego.NSRouter("/blog.comment.get", &common.Controller{}, "post:BlogCommonGet"),
			beego.NSRouter("/blog.view", &common.Controller{}, "post:BlogView"),
			beego.NSRouter("/talk.get", &common.Controller{}, "post:TalkGet"),
			beego.NSRouter("/talk.comment.get", &common.Controller{}, "post:TalkCommentGet"),
			beego.NSRouter("/talk.like", &common.Controller{}, "post:TalkLike"),
			beego.NSRouter("/blog.comment.like", &common.Controller{}, "post:BlogCommonLike"),
			beego.NSRouter("/profile.blog.get", &common.Controller{}, "post:ProfileBlogGet"),
			beego.NSRouter("/profile.talk.get", &common.Controller{}, "post:ProfileTalkGet"),
			beego.NSRouter("/profile.collect.get", &common.Controller{}, "post:ProfileCollectGet"),
			beego.NSRouter("/user.profile.get", &common.Controller{}, "post:UserProfileGet"),
			beego.NSRouter("/collect.profile.get", &common.Controller{}, "post:CollectProfileGet"),
			beego.NSRouter("/user.profile.care", &common.Controller{}, "post:UserProfileCare"),
			beego.NSRouter("/user.profile.cared", &common.Controller{}, "post:UserProfileCared"),
			beego.NSRouter("/profile.leave.get", &common.Controller{}, "post:ProfileLeaveGet"),
			beego.NSRouter("/blog.search", &common.Controller{}, "post:BlogSearch"),
			beego.NSRouter("/talk.search", &common.Controller{}, "post:TalkSearch"),
			beego.NSRouter("/user.search", &common.Controller{}, "post:UserSearch"),
		),
	)
	beego.AddNamespace(ns)
}

func init() {
	apiRegister()
	log.Info("router init successfully")
}
