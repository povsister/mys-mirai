package post

import (
	"strconv"

	"github.com/povsister/mys-mirai/mys/api/meta"
	"github.com/povsister/mys-mirai/mys/api/response/post"
	"github.com/povsister/mys-mirai/mys/rest"
)

type PostInterface interface {
	Get(pid int, opt meta.GetPostOptions) (*post.FullPostInfoResponse, error)
}

type postImpl struct {
	client rest.Interface
	gid    rest.GameType
}

func newPostImpl(c rest.Interface, gid rest.GameType) *postImpl {
	return &postImpl{client: c, gid: gid}
}

func (c *postImpl) Get(pid int, opt meta.GetPostOptions) (ret *post.FullPostInfoResponse, err error) {
	ret = &post.FullPostInfoResponse{}
	err = c.client.Get().
		Use(opt).
		GID(c.gid).
		Path("/post/wapi/getPostFull").
		ParamSet("post_id", strconv.Itoa(pid)).
		ParamSet("read", "1").
		Do().Into(ret)
	return
}
