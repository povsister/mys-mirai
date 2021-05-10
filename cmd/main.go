package main

import (
	"github.com/povsister/mys-mirai/bot"
	"github.com/povsister/mys-mirai/config"
	"github.com/povsister/mys-mirai/pkg/log"
	"github.com/povsister/mys-mirai/pkg/util"
)

func main() {
	cfg := config.Read()
	if cfg == nil {
		return
	}
	b := bot.NewBot(cfg.Account.Uin, cfg.Account.Password)
	if err := b.Login(); err != nil {
		log.Error().Err(err).Msg("Bot启动失败")
		return
	}

	log.Info().Msg("Bot启动完成 Ctrl+C 退出")
	// wait for interrupt signal
	<-util.SetupSignalHandler()
}
