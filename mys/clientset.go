package mys

import (
	"net/url"

	"github.com/povsister/mys-mirai/mys/api/request/moderator"
	"github.com/povsister/mys-mirai/mys/api/request/user"
	"github.com/povsister/mys-mirai/mys/rest"
)

var MysApiBase = url.URL{
	Host:   "bbs-api.mihoyo.com",
	Scheme: "https",
	Path:   "/",
}

type Clientset struct {
	restClient rest.Interface
	moderator  *moderator.ModeratorClient
	user       *user.UserClient
}

func (c *Clientset) RESTClient() rest.Interface {
	return c.restClient
}

func (c *Clientset) Moderator() *moderator.ModeratorClient {
	return c.moderator
}

func (c *Clientset) User() *user.UserClient {
	return c.user
}

func NewClient(config *rest.Config) *Clientset {
	restClient := rest.NewRESTClient(MysApiBase, config)

	return &Clientset{
		restClient: restClient,
		moderator:  moderator.NewModeratorClient(restClient),
		user:       user.NewUserClient(restClient),
	}
}
