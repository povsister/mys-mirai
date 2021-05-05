package bot

import "github.com/Mrs4s/MiraiGo/client"

type Bot struct {
	c *client.QQClient
}

func NewBot(uid int64, pw string) *Bot {
	b := &Bot{
		c: client.NewClient(uid, pw),
	}
	return b
}
