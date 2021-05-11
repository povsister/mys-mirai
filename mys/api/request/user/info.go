package user

import (
	"github.com/povsister/mys-mirai/mys/api/request/meta"
	"github.com/povsister/mys-mirai/mys/api/response/user"
	"github.com/povsister/mys-mirai/mys/rest"
	"strconv"
)

type UserInfoInterface interface {
	Get(uid int, opt meta.UserInfoGetOptions) (*user.FullUserInfo, error)
}

type userInfoManager struct {
	client rest.Interface
	gid    rest.GameType
}

func newUserInfoManager(c rest.Interface, gid rest.GameType) *userInfoManager {
	return &userInfoManager{client: c, gid: gid}
}

func (c *userInfoManager) Get(uid int, opt meta.UserInfoGetOptions) (ret *user.FullUserInfo, err error) {
	ret = &user.FullUserInfo{}
	req := c.client.Get().
		Use(opt).
		GID(c.gid).
		Path("/user/wapi/getUserFullInfo")
	if uid != meta.UserMyself {
		req.ParamSet("uid", strconv.Itoa(uid))
	}
	err = req.Do().Into(ret)
	return
}
