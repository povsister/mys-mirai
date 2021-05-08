package moderator

import (
	"github.com/povsister/mys-mirai/mys/api/request/meta"
	"github.com/povsister/mys-mirai/mys/rest"
)

type PostInterface interface {
	Delete(pid int, opt meta.DeletePostOptions) error
}

type postManger struct {
	client rest.Interface
	gid    meta.GameType
}

func newPostManger(c *ModeratorClient, forum meta.GameType) *postManger {
	return &postManger{client: c.restClient, gid: forum}
}

func (c *postManger) Delete(pid int, opt meta.DeletePostOptions) error {
	r := c.client.Post().
		GID(c.gid).
		Path("/post/wapi/deletePost").
		BodyKV("post_id", pid)
	if opt.Reason != meta.NoReason {
		r.BodyKV("config_id", opt.Reason)
	}
	return r.Do().Error()
}

func (c *postManger) Move(pid int, opt meta.MovePostOptions) error {
	return c.client.Post().
		GID(c.gid).
		Path("/post/wapi/movePost").
		BodyKV("post_id", pid).
		BodyKV("f_forum_id", opt.To).
		Do().Error()
}
