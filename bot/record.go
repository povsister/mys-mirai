package bot

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/message"
)

func (b *Bot) RecordMessage(m interface{}) {
	if b.msgCache == nil {
		return
	}
	DispatchMessage(m, DispatchTo{
		From: "bot.RecordMessage",
		Group: func(gm *message.GroupMessage) {
			b.recordGroupMsg(gm)
		},
		Private: func(pm *message.PrivateMessage) {
			b.recordPrivateMsg(pm)
		},
	})
}

func (b *Bot) recordGroupMsg(m *message.GroupMessage) {
	b.msgCache.Add(toGroupMsgKey(m.GroupCode, m.Id), m)
}

func (b *Bot) recordPrivateMsg(m *message.PrivateMessage) {
	b.msgCache.Add(toPrivateMsgKey(m.Sender.Uin, m.Id), m)
}

func toGroupMsgKey(gid int64, mid int32) string {
	return fmt.Sprintf("g.%d.%d", gid, mid)
}

func toPrivateMsgKey(sender int64, mid int32) string {
	return fmt.Sprintf("p.%d.%d", sender, mid)
}

func (b *Bot) LookupGroupMessage(groupCode int64, messageID int32) *message.GroupMessage {
	if b.msgCache == nil {
		return nil
	}
	m, ok := b.msgCache.Get(toGroupMsgKey(groupCode, messageID))
	if !ok {
		return nil
	}
	return m.(*message.GroupMessage)
}

func (b *Bot) LookupPrivateMessage(sender int64, messageID int32) *message.PrivateMessage {
	if b.msgCache == nil {
		return nil
	}
	m, ok := b.msgCache.Get(toPrivateMsgKey(sender, messageID))
	if !ok {
		return nil
	}
	return m.(*message.PrivateMessage)
}
