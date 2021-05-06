package event

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

var PrivateMessageListener = new(privateMessage)

type privateMessage struct {
}

func (m *privateMessage) Register(c *client.QQClient) {
	c.OnPrivateMessage(func(cli *client.QQClient, msg *message.PrivateMessage) {
		msgStr := msg.ToString()
		if len(msgStr) == 0 {
			return
		}
		reply := message.NewSendingMessage().Append(message.NewText("喵~ " + msgStr))
		cli.SendPrivateMessage(msg.Sender.Uin, reply)
	})
}
