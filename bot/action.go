package bot

import (
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/povsister/mys-mirai/pkg/qqmsg"
)

func (b *Bot) ReplyStr(m interface{}, reply string) {
	switch msg := m.(type) {
	case *message.GroupMessage:
		b.replyGroupMsg(msg, reply)
	case *message.PrivateMessage:
		b.replyPrivateMsg(msg, reply)
	default:
		logger.Error().Str("msg", reply).Msgf("can not reply message. unknown message type %T", m)
	}
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
	switch msg := m.(type) {
	case *message.GroupMessage:
		b.sendGroupMsg(msg.GroupCode, toSend)
	case *message.PrivateMessage:
		b.sendPrivateMsg(msg.Sender.Uin, toSend)
	default:
		logger.Error().Str("msg", toSend).Msgf("can not send message by msg. unknown message type %T", m)
	}
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
	if qqmsg.IsPrivateMsg(m) {
		b.SendStrByMsg(m, toSend)
		return
	}
	b.ReplyStr(m, toSend)
}
