package moderator

import (
	"github.com/povsister/mys-mirai/mys/api/request/meta"
	"github.com/povsister/mys-mirai/mys/rest"
)

type PostInterface interface {
	Delete(pid int, opt meta.DeletePostOptions) error
	Move(pid int, opt meta.MovePostOptions) error
	RemoveTopic(pid int, opt meta.RemoveTopicOptions) error
}

type postManger struct {
	client rest.Interface
	gid    rest.GameType
}

func newPostManger(c *ModeratorClient, gid rest.GameType) *postManger {
	return &postManger{client: c.restClient, gid: gid}
}

func (c *postManger) Delete(pid int, opt meta.DeletePostOptions) error {
	return c.client.Post().
		Use(opt).
		GID(c.gid).
		Path("/post/wapi/deletePost").
		BodyKV("post_id", pid).
		Do().Error()
}

func (c *postManger) Move(pid int, opt meta.MovePostOptions) error {
	return c.client.Post().
		Use(opt).
		GID(c.gid).
		Path("/post/wapi/movePost").
		BodyKV("post_id", pid).
		Do().Error()
}

func (c *postManger) RemoveTopic(pid int, opt meta.RemoveTopicOptions) error {
	return c.client.Post().
		Use(opt).
		GID(c.gid).
		Path("/post/wapi/removePostTopicsByTopicIDs").
		BodyKV("post_id", pid).
		Do().Error()
}
