package mys

import (
	"github.com/povsister/mys-mirai/mys/api/request/moderator"
	"github.com/povsister/mys-mirai/mys/rest"
	"net/url"
)

var MysApiBase = url.URL{
	Host:   "bbs-api.mihoyo.com",
	Scheme: "https",
	Path:   "/",
}

type Clientset struct {
	moderator *moderator.ModeratorClient
}

func (c *Clientset) Moderator() *moderator.ModeratorClient {
	return c.moderator
}

func NewClient(config *rest.Config) *Clientset {
	restClient := rest.NewRESTClient(MysApiBase, config)

	return &Clientset{
		moderator: moderator.NewModeratorClient(restClient),
	}
}
