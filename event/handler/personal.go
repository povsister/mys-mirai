package handler

import (
	"bytes"
	"github.com/povsister/mys-mirai/bot"
	"github.com/povsister/mys-mirai/event/router"
	"github.com/povsister/mys-mirai/mys/api/request/meta"
	"github.com/povsister/mys-mirai/mys/rest"
)

var personalHandlers = map[string]handlerDetail{
	"/个人信息": {"查看关联的米游社帐号", handlePersonalInfo},
}

func init() {
	for path, h := range personalHandlers {
		router.Router.RegisterHandler(router.Message, path, h.Handler)
	}
	router.Router.RegisterHandler(router.Message, "/帮助", handleHelp)
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

func handleHelp(b *bot.Bot, m interface{}, _ router.Params) {
	buf := new(bytes.Buffer)
	buf.WriteString("喵~\n")
	for cmd, h := range moderatorHandlers {
		buf.WriteString(cmd + "  " + h.Descr + "\n")
	}
	buf.WriteString("\n")
	for cmd, h := range personalHandlers {
		buf.WriteString(cmd + "  " + h.Descr + "\n")
	}

	b.SendStrByMsg(m, buf.String())
}
