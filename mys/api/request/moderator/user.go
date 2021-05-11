package moderator

import (
	"github.com/povsister/mys-mirai/mys/api/request/meta"
	"github.com/povsister/mys-mirai/mys/rest"
)

type UserMngInterface interface {
	Silence(uid int, opt meta.SilenceUserOptions) error
	UnSilence(uid int, opt meta.UnSilenceUserOptions) error
}

type userManager struct {
	client rest.Interface
	gid    rest.GameType
}

func newUserManager(c *ModeratorClient, gid rest.GameType) *userManager {
	return &userManager{client: c.restClient, gid: gid}
}

func (c *userManager) Silence(uid int, opt meta.SilenceUserOptions) error {
	return c.client.Post().
		Use(opt).
		GID(c.gid).
		Path("/user/wapi/silenceUser").
		BodyKV("uid", uid).
		Do().Error()
}

func (c *userManager) UnSilence(uid int, opt meta.UnSilenceUserOptions) error {
	return c.client.Post().
		Use(opt).
		GID(c.gid).
		Path("/user/wapi/unSilenceUser").
		BodyKV("uid", uid).
		Do().Error()
}
