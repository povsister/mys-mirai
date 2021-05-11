package event

import (
	"bytes"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/povsister/mys-mirai/bot"
	"github.com/povsister/mys-mirai/event/basic"
	"github.com/povsister/mys-mirai/mys/api/request/meta"
	"github.com/povsister/mys-mirai/mys/rest"
	"strings"
	"text/template"
)

var PrivateMessageListener = &privateMessage{new(basic.BasicListener)}

func init() {
	bot.RegisterEvent(PrivateMessageListener)
}

type privateMessage struct {
	*basic.BasicListener
}

func (pm *privateMessage) Start() {
	pm.Bot.Q().OnPrivateMessage(func(c *client.QQClient, m *message.PrivateMessage) {
		msg := m.ToString()
		if len(msg) == 0 {
			return
		}
		if strings.HasPrefix(msg, "/个人信息") {
			mys := pm.Bot.M().Get(m.Sender.Uin)
			if mys == nil {
				reply := message.NewSendingMessage().Append(message.NewText("尚未登录"))
				c.SendPrivateMessage(m.Sender.Uin, reply)
				return
			}
			info, err := mys.User().Info(rest.Genshin).Get(meta.UserMyself, meta.UserInfoGetOptions{})
			if err != nil {
				reply := message.NewSendingMessage().Append(message.NewText("错误: " + err.Error()))
				c.SendPrivateMessage(m.Sender.Uin, reply)
				return
			}
			if !info.LoggedIn() {
				reply := message.NewSendingMessage().Append(message.NewText("错误: " + info.Error()))
				c.SendPrivateMessage(m.Sender.Uin, reply)
				return
			}
			text, err := template.New("personal").Parse(personal)
			if err != nil {
				reply := message.NewSendingMessage().Append(message.NewText("错误: " + err.Error()))
				c.SendPrivateMessage(m.Sender.Uin, reply)
				return
			}
			buf := new(bytes.Buffer)
			_ = text.Execute(buf, info)
			reply := message.NewSendingMessage().Append(message.NewText(buf.String()))
			c.SendPrivateMessage(m.Sender.Uin, reply)
		} else {
			reply := message.NewSendingMessage().Append(message.NewText("喵~ " + m.Sender.DisplayName()))
			c.SendPrivateMessage(m.Sender.Uin, reply)
		}
	})
}

var personal = `已关联的米游社用户信息
昵称: {{.UserInfo.Nickname}}
通行证: {{.UserInfo.UID}}`
