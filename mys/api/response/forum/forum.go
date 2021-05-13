package forum

import (
	"github.com/povsister/mys-mirai/mys/api/meta"
	"github.com/povsister/mys-mirai/mys/rest"
	"github.com/povsister/mys-mirai/mys/runtime"
)

type (
	// /api/community/forum/home/forums?gids=2
	ForumListResponse struct {
		runtime.ObjectMeta `json:",inline"`
		ForumListData      `json:"data"`
	}

	ForumListData struct {
		ForumList []ForumInfoDetail `json:"forumlists"`
	}
)
type (
	ForumInfoDetail struct {
		ForumInfoBasic `json:",inline"`
		CreatedAt      runtime.TimeRFC3339 `json:"created_at"`
		UpdatedAt      runtime.TimeRFC3339 `json:"updated_at"`
		HeaderImgURL   string              `json:"header_image"`
		IconPureURL    string              `json:"icon_pure"`
		Description    string              `json:"des"`
		MaxTop         runtime.Int         `json:"max_top"`
		// TODO: finish remaining fields
	}

	ForumInfoBasic struct {
		ID      meta.Forum    `json:"id"`
		GID     rest.GameType `json:"game_id"`
		IconURL string        `json:"icon"`
		Name    string        `json:"name"`
	}
)
