package user

import "github.com/povsister/mys-mirai/mys/runtime"

type (
	// /collection/wapi/collection/userList?gids=2&size=20&uid=
	CollectionListResponse struct {
		runtime.ObjectMeta `json:",inline"`
		CollectionListData `json:"data"`
	}

	CollectionListData struct {
		IsLast     bool               `json:"is_last"`
		NextOffset int                `json:"next_offset"`
		List       []CollectionDetail `json:"list"`
	}

	CollectionDetail struct {
		ID            runtime.Int           `json:"id"`
		Title         string                `json:"title"`
		CoverURL      string                `json:"cover"`
		Description   string                `json:"desc"`
		CreatorUID    runtime.Int           `json:"uid"`
		IsDeleted     bool                  `json:"is_deleted"`
		IsFollowing   bool                  `json:"is_following"`
		PostUpdatedAt runtime.UnixTimestamp `json:"post_updated_at"`
		Total         int                   `json:"post_num"`
		ViewNum       int                   `json:"view_num"`
	}
)
