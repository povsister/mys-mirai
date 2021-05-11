package handler

import (
	"github.com/povsister/mys-mirai/bot"
	"github.com/povsister/mys-mirai/event/router"
	"github.com/povsister/mys-mirai/mys/api/request/meta"
	"github.com/povsister/mys-mirai/mys/rest"
)

func init() {
	router.Router.RegisterHandler(router.Message, "/个人信息", handlePersonalInfo)
}

func handlePersonalInfo(b *bot.Bot, m interface{}, _ router.Params) {
	mys := b.MByMsg(m)
	if mys == nil {
		return
	}
	info, err := mys.User().Info(rest.Genshin).Get(meta.UserMyself, meta.UserInfoGetOptions{})
	if err != nil {
		b.ReplyStr(m, runtimeErr(err))
		return
	}
	if !info.LoggedIn() {
		b.ReplyStr(m, "契约出现问题了喵~\n"+info.Error())
		return
	}
	b.SendOrReplyStrByMsg(m, useTemplate(personalInfo, info))
}
