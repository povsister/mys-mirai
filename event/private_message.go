package event

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/povsister/mys-mirai/bot"
	"github.com/povsister/mys-mirai/event/basic"
	_ "github.com/povsister/mys-mirai/event/handler"
	"github.com/povsister/mys-mirai/event/router"
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
		router.Router.HandleMessage(pm.Bot, m)
	})
}
