package bot

import (
	"github.com/Mrs4s/MiraiGo/message"
)

func (b *Bot) ReplyStr(m interface{}, reply string) {
	DispatchMessage(m, DispatchTo{
		From: "bot.ReplyStr",
		Group: func(gm *message.GroupMessage) {
			b.replyGroupMsg(gm, reply)
		},
		Private: func(pm *message.PrivateMessage) {
			b.replyPrivateMsg(pm, reply)
		},
	})
}

func (b *Bot) replyGroupMsg(m *message.GroupMessage, reply string) {
	msg := message.NewSendingMessage()
	msg.Append(message.NewReply(m))
	msg.Append(message.NewText(reply))
	b.c.SendGroupMessage(m.GroupCode, msg)
}

func (b *Bot) replyPrivateMsg(m *message.PrivateMessage, reply string) {
	msg := message.NewSendingMessage()
	msg.Append(message.NewPrivateReply(m))
	msg.Append(message.NewText(reply))
	b.c.SendPrivateMessage(m.Sender.Uin, msg)
}

func (b *Bot) SendStrByMsg(m interface{}, toSend string) {
	DispatchMessage(m, DispatchTo{
		From: "bot.SendStrByMsg",
		Group: func(gm *message.GroupMessage) {
			b.sendGroupMsg(gm.GroupCode, toSend)
		},
		Private: func(pm *message.PrivateMessage) {
			b.sendPrivateMsg(pm.Sender.Uin, toSend)
		},
	})
}

func (b *Bot) sendGroupMsg(groupCode int64, send string) {
	msg := message.NewSendingMessage()
	msg.Append(message.NewText(send))
	b.c.SendGroupMessage(groupCode, msg)
}

func (b *Bot) sendPrivateMsg(uin int64, send string) {
	msg := message.NewSendingMessage()
	msg.Append(message.NewText(send))
	b.c.SendPrivateMessage(uin, msg)
}

func (b *Bot) SendOrReplyStrByMsg(m interface{}, toSend string) {
	if IsPrivateMsg(m) {
		b.SendStrByMsg(m, toSend)
		return
	}
	b.ReplyStr(m, toSend)
}
