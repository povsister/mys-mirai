package post

import "github.com/povsister/mys-mirai/mys/runtime"

type (
	// post/wapi/getForumPostList?forum_id=29&gids=2&is_good=false&is_hot=false&page_size=20&sort_type=1
	ForumPostListResponse struct {
		runtime.ObjectMeta `json:",inline"`
		ForumPostListData  `json:"data"`
	}

	ForumPostListData struct {
		IsLast   bool                  `json:"is_last"`
		IsOrigin bool                  `json:"is_origin"`
		LastID   runtime.UnixTimestamp `json:"last_id"`
		List     []PostInfo            `json:"list"`
	}
)
