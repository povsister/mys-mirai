package post

import "github.com/povsister/mys-mirai/mys/rest"

type PostClient struct {
	restClient rest.Interface
}

func NewPostClient(c rest.Interface) *PostClient {
	return &PostClient{restClient: c}
}

func (c *PostClient) Info(gid rest.GameType) PostInterface {
	return newPostImpl(c.restClient, gid)
}
