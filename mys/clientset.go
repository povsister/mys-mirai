package mys

import "github.com/povsister/mys-mirai/mys/api/request/moderator"

type Clientset struct {
	moderator *moderator.ModeratorClient
}

func (c *Clientset) Moderator() *moderator.ModeratorClient {
	return c.moderator
}
