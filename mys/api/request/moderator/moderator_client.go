package moderator

import (
	"github.com/povsister/mys-mirai/mys/api/request/meta"
	"github.com/povsister/mys-mirai/mys/rest"
)

type ModeratorClient struct {
	restClient rest.Interface
}

func (c *ModeratorClient) Post(forum meta.GameType) PostInterface {
	return newPostManger(c, forum)
}
