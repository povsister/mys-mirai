package user

import "github.com/povsister/mys-mirai/mys/rest"

type UserClient struct {
	restClient rest.Interface
}

func NewUserClient(c rest.Interface) *UserClient {
	return &UserClient{restClient: c}
}

func (c *UserClient) Info(gid rest.GameType) UserInfoInterface {
	return newUserInfoManager(c.restClient, gid)
}
