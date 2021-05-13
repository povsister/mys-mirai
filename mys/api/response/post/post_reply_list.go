package post

import (
	"github.com/povsister/mys-mirai/mys/api/meta"
	"github.com/povsister/mys-mirai/mys/api/response/user"
	"github.com/povsister/mys-mirai/mys/rest"
	"github.com/povsister/mys-mirai/mys/runtime"
)

type (
	// /post/wapi/getPostReplies?gids=2&is_hot=true&post_id=xxx&size=20
	ReplyListResponse struct {
		runtime.ObjectMeta `json:",inline"`
		ReplyListData      `json:"data"`
	}

	ReplyListData struct {
		IsLast       bool        `json:"is_last"`
		LastID       runtime.Int `json:"last_id"`
		PostAuthorID runtime.Int `json:"post_owner_id"`
		List         []ReplyInfo `json:"list"`
	}
)

type (
	// TODO: recursive reply type
	ReplyInfo struct {
		User         user.UserInfoBasic `json:"user"`
		IsPostAuthor bool               `json:"is_lz"`
		ReplyToPost  PostInfo           `json:"r_post"`
		//ReplyToReply  ReplyInfo          `json:"r_reply"`
		ReplyToUser  user.UserInfoBasic `json:"r_user"`
		Stat         ReplyStat          `json:"stat"`
		MasterStatus MasterStatus       `json:"master_status"`
		//SubReply      []ReplyInfo        `json:"sub_replies"`
		SubReplyCount int           `json:"sub_reply_count"`
		MyOperation   SelfOperation `json:"self_operation"`
		Detail        ReplyContent  `json:"reply"`
	}

	ReplyStat struct {
		Likes      int `json:"like_num"`
		Replies    int `json:"reply_num"`
		SubReplies int `json:"sub_num"`
	}
)

type (
	ReplyContent struct {
		FloorID        runtime.Int           `json:"floor_id"`
		ReplyID        runtime.Int           `json:"reply_id"`
		PostID         runtime.Int           `json:"post_id"`
		AuthorID       runtime.Int           `json:"uid"`
		ReplyToUID     runtime.Int           `json:"r_uid"`
		BodyHTML       string                `json:"content"`
		BodyStructured string                `json:"structured_content"`
		Forum          meta.Forum            `json:"f_forum_id"`
		GID            rest.GameType         `json:"game_id"`
		ReplyToReplyID runtime.Int           `json:"f_reply_id"`
		HasBlockedWord bool                  `json:"has_blocked_word"`
		IsTop          bool                  `json:"is_top"`
		IsDeleted      DeleteStatus          `json:"is_deleted"`
		DeleteSrc      DeleteSrcType         `json:"delete_src"`
		CreatedAt      runtime.UnixTimestamp `json:"created_at"`
		DeletedAt      runtime.UnixTimestamp `json:"deleted_at"`
		UpdatedAt      runtime.UnixTimestamp `json:"updated_at"`
	}

	DeleteSrcType uint8
)
