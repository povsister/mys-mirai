package event

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/povsister/mys-mirai/bot"
	"github.com/povsister/mys-mirai/event/basic"
	_ "github.com/povsister/mys-mirai/event/handler"
	"github.com/povsister/mys-mirai/event/router"
)

var GroupMessageListener = &groupMessage{new(basic.BasicListener)}

func init() {
	bot.RegisterEvent(GroupMessageListener)
}

type groupMessage struct {
	*basic.BasicListener
}

func (gm *groupMessage) Start() {
	gm.Bot.Q().OnGroupMessage(func(_ *client.QQClient, msg *message.GroupMessage) {
		router.Router.HandleMessage(gm.Bot, msg)
	})
}
