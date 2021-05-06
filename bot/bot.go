package bot

import (
	"fmt"
	"sync"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/povsister/mys-mirai/pkg/log"
	"github.com/povsister/mys-mirai/pkg/util/fs"
)

type EventListener interface {
	Register(client *client.QQClient)
}

type Bot struct {
	c  *client.QQClient
	lm *loginManger
}

// loginManager 负责登录以及断线重连相关的活
type loginManger struct {
	bot *Bot
	// 是否已经登录 用来区分重连
	alreadyLogin bool
	// 重连次数限制
	reConnLimit int
	// 保证同时只有一次登录或者断线重连
	mu   sync.Mutex
	once sync.Once
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

func (b *Bot) RegisterEvent(l EventListener) {
	l.Register(b.c)
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
