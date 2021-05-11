package mys

import (
	"github.com/povsister/mys-mirai/mys/api/request/post"
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
	qid        int64
	restClient rest.Interface
	moderator  *moderator.ModeratorClient
	user       *user.UserClient
	post       *post.PostClient
}

func (c *Clientset) Qid() int64 {
	return c.qid
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

func (c *Clientset) Post() *post.PostClient {
	return c.post
}

func NewClient(config *rest.Config) *Clientset {
	restClient := rest.NewRESTClient(MysApiBase, config)

	return &Clientset{
		qid:        config.Qid,
		restClient: restClient,
		moderator:  moderator.NewModeratorClient(restClient),
		user:       user.NewUserClient(restClient),
		post:       post.NewPostClient(restClient),
	}
}
