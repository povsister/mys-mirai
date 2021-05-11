package bot

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	lru "github.com/hashicorp/golang-lru"
	"github.com/povsister/mys-mirai/mys"
	"github.com/povsister/mys-mirai/pkg/log"
	"github.com/povsister/mys-mirai/pkg/util/fs"
)

var logger log.Logger

func init() {
	logger = log.SubLogger("bot")
}

// 由 event listener 实现
type EventListener interface {
	Register(c *Bot)
	Start()
}

// 所有注册过的event listener
var registeredListener = make([]EventListener, 0, 10)

// 用于event向bot注册event listener
func RegisterEvent(e EventListener) {
	registeredListener = append(registeredListener, e)
}

// bot
type Bot struct {
	c        *client.QQClient
	lm       *loginManger
	mys      *mys.UserManager
	msgCache *lru.ARCCache
}

func NewBot(uid int64, pw string) *Bot {
	b := &Bot{}
	b.lm = &loginManger{bot: b, reConnLimit: 10}
	b.loadDeviceJSON()
	b.c = client.NewClient(uid, pw)
	b.setupLogHandler()
	b.handleServerUpdated()
	b.mys = mys.NewUserManager()
	lruc, err := lru.NewARC(5000)
	if err != nil {
		log.Warn().Err(err).Msg("缓存初始化失败 将无法正常记录历史消息")
	} else {
		b.msgCache = lruc
	}
	return b
}

func (b *Bot) Q() *client.QQClient {
	return b.c
}

func (b *Bot) M() *mys.UserManager {
	return b.mys
}

// 从消息自动获取对应的 mys client
func (b *Bot) MByMsg(m interface{}) *mys.Clientset {
	uin := ExtractMsgSenderUin(m)
	if uin == 0 {
		return nil
	}
	ret := b.mys.Get(uin)
	if ret == nil {
		b.ReplyStr(m, "契约尚未签订喵~")
		return nil
	}
	return ret
}

// 启动所有注册过的 event listener
func (b *Bot) registerKnownEvents() {
	for _, e := range registeredListener {
		e.Register(b)
		e.Start()
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
