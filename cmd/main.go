package main

import (
	"github.com/povsister/mys-mirai/bot"
	"github.com/povsister/mys-mirai/config"
	"github.com/povsister/mys-mirai/event"
	"github.com/povsister/mys-mirai/pkg/log"
	"github.com/povsister/mys-mirai/pkg/util"
)

func main() {
	log.InitLogger()

	cfg := config.Read()
	if cfg == nil {
		return
	}
	b := bot.NewBot(cfg.Account.Uin, cfg.Account.Password)
	if err := b.Login(); err != nil {
		log.Error().Err(err).Msg("Bot启动失败")
		return
	}

	b.RegisterEvent(event.GroupMessageListener)
	b.RegisterEvent(event.PrivateMessageListener)

	// wait for interrupt signal
	<-util.SetupSignalHandler()
}
