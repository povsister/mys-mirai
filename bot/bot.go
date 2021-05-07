package bot

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/povsister/mys-mirai/pkg/log"
	"github.com/povsister/mys-mirai/pkg/util/fs"
)

// 由 event listener 实现
type EventListener interface {
	Register(client *client.QQClient)
}

// 所有注册过的event listener
var registeredListener = make([]EventListener, 0, 10)

// 用于event向bot注册event listener
func RegisterEvent(e EventListener) {
	registeredListener = append(registeredListener, e)
}

// bot
type Bot struct {
	c  *client.QQClient
	lm *loginManger
}

func NewBot(uid int64, pw string) *Bot {
	b := &Bot{}
	b.lm = &loginManger{bot: b, reConnLimit: 10}
	b.loadDeviceJSON()
	b.c = client.NewClient(uid, pw)
	b.setupLogHandler()
	b.handleServerUpdated()
	return b
}

// 启动所有注册过的 event listener
func (b *Bot) registerKnownEvents() {
	for _, e := range registeredListener {
		e.Register(b.c)
	}
}

func (b *Bot) handleServerUpdated() {
	b.c.OnServerUpdated(func(c *client.QQClient, e *client.ServerUpdatedEvent) bool {
		// 默认接受服务端地址更新
		return true
	})
}

func (b *Bot) sessionFile() string {
	if b.c.Uin != 0 {
		return fmt.Sprintf("session.%d.token", b.c.Uin)
	}
	return "session.token"
}

func (b *Bot) saveSessionFile(data []byte) error {
	sf := b.sessionFile()
	err := fs.WriteFile(sf, data, 0644)
	if err == nil {
		log.Info().Msgf("session.token 已保存至 %s", sf)
	}
	return err
}
