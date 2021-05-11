package basic

import "github.com/povsister/mys-mirai/bot"

type BasicListener struct {
	Bot *bot.Bot
}

func (l *BasicListener) Register(b *bot.Bot) {
	l.Bot = b
}
