package event

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/povsister/mys-mirai/bot"
	"strings"
)

var GroupMessageListener = new(groupMessage)

func init() {
	bot.RegisterEvent(GroupMessageListener)
}

type groupMessage struct {
}

func (m *groupMessage) Register(c *client.QQClient) {
	c.OnGroupMessage(func(cli *client.QQClient, msg *message.GroupMessage) {
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
