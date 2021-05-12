package forum

import (
	"github.com/povsister/mys-mirai/mys/api/meta"
	"github.com/povsister/mys-mirai/mys/rest"
	"github.com/povsister/mys-mirai/mys/runtime"
)

type ForumInfoDetail struct {
	ForumInfoBasic `json:",inline"`
}

type ForumInfoBasic struct {
	ID          meta.Forum    `json:"id"`
	Gid         rest.GameType `json:"game_id"`
	IconURL     string        `json:"icon"`
	IconPureURL string        `json:"icon_pure"`
	Name        string        `json:"name"`
	Descr       string        `json:"des"`
	MaxTop      runtime.Int   `json:"max_top"`
}
