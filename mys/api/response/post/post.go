package post

import (
	"github.com/povsister/mys-mirai/mys/api/response/user"
	"github.com/povsister/mys-mirai/mys/runtime"
)

type FullPostInfo struct {
	runtime.ObjectMeta `json:",inline"`
	FullPostInfoData   `json:"data"`
}

type FullPostInfoData struct {
	Post PostInfo `json:"post"`
}

type PostInfo struct {
	// TODO: finish all fields
	User user.UserInfoBasic `json:"user"`
}
