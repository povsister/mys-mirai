package moderator

import (
	"github.com/povsister/mys-mirai/mys/rest"
)

type ModeratorClient struct {
	restClient rest.Interface
}

func NewModeratorClient(c rest.Interface) *ModeratorClient {
	return &ModeratorClient{restClient: c}
}

func (c *ModeratorClient) Post(gid rest.GameType) PostMngInterface {
	return newPostManger(c, gid)
}

func (c *ModeratorClient) User(gid rest.GameType) UserMngInterface {
	return newUserManager(c, gid)
}
