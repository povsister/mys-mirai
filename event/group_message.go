package event

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/povsister/mys-mirai/bot"
	"github.com/povsister/mys-mirai/event/basic"
	"strings"
)

var GroupMessageListener = &groupMessage{new(basic.BasicListener)}

func init() {
	bot.RegisterEvent(GroupMessageListener)
}

type groupMessage struct {
	*basic.BasicListener
}

func (gm *groupMessage) Start() {
	gm.Bot.Q().OnGroupMessage(func(cli *client.QQClient, msg *message.GroupMessage) {
		msgStr := msg.ToString()
		if len(msgStr) == 0 {
			return
		}
		if strings.Contains(msgStr, "Hello") {
			msg.Sender.DisplayName()
			reply := message.NewSendingMessage().Append(message.NewText("Hello " + msg.Sender.DisplayName()))
			cli.SendGroupMessage(msg.GroupCode, reply)
		}
	})
}
