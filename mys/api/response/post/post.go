package post

import (
	"github.com/povsister/mys-mirai/mys/api/meta"
	"github.com/povsister/mys-mirai/mys/api/response/common"
	"github.com/povsister/mys-mirai/mys/api/response/forum"
	"github.com/povsister/mys-mirai/mys/api/response/user"
	"github.com/povsister/mys-mirai/mys/rest"
	"github.com/povsister/mys-mirai/mys/runtime"
)

type (
	// /post/wapi/getPostFull?gids=2&post_id=xxx&read=1
	FullPostInfoResponse struct {
		runtime.ObjectMeta `json:",inline"`
		FullPostInfoData   `json:"data"`
	}

	FullPostInfoData struct {
		Post PostInfo `json:"post"`
	}
)

type (
	PostInfo struct {
		Collection       PostCollection        `json:"collection"`
		Cover            Cover                 `json:"cover"`
		Forum            forum.ForumInfoBasic  `json:"forum"`
		User             user.UserInfoBasic    `json:"user"`
		MyOperation      SelfOperation         `json:"self_operation"`
		Stat             PostStat              `json:"stat"`
		TopicList        []TopicDetail         `json:"topics"`
		Detail           PostContent           `json:"post"`
		ImageList        []Image               `json:"image_list"`
		VodList          []VodObject           `json:"vod_list"`
		VoteCount        int                   `json:"vote_count"`
		IsBlockOn        bool                  `json:"is_block_on"`
		IsOfficialMaster bool                  `json:"is_official_master"`
		IsUserMaster     bool                  `json:"is_user_master"`
		HasHotReplies    bool                  `json:"has_hot_replies"`
		LastModifiedAt   runtime.UnixTimestamp `json:"last_modify_time"`
	}

	PostStat struct {
		Bookmarks int `json:"bookmark_num"`
		Likes     int `json:"like_num"`
		Replies   int `json:"reply_num"`
		Views     int `json:"view_num"`
	}
)

type (
	PostCollection struct {
		ID          runtime.Int   `json:"collection_id"`
		Title       string        `json:"collection_title"`
		CurrentPos  runtime.Int   `json:"cur"`
		Total       runtime.Int   `json:"total"`
		NextPostID  runtime.Int   `json:"next_post_id"`
		NextPostGID rest.GameType `json:"next_post_game_id"`
		PrevPostID  runtime.Int   `json:"prev_post_id"`
		PrevPostGID rest.GameType `json:"prev_post_game_id"`
	}
)

type (
	Cover struct {
		URL              string      `json:"url"`
		Crop             common.Crop `json:"crop"`
		Format           string      `json:"format"`
		SetByUser        bool        `json:"is_user_set_cover"`
		common.SizeByte  `json:",inline"`
		common.Rectangle `json:",inline"`
	}
)

type (
	SelfOperation struct {
		Attitude   SelfAttitude `json:"attitude"`
		Bookmarked bool         `json:"is_collected"`
	}

	SelfAttitude int
)

const (
	// 无点赞
	Noop SelfAttitude = 0
	// 点赞
	Like SelfAttitude = 1
)

type (
	Image struct {
		URL    string      `json:"url"`
		Crop   common.Crop `json:"crop"`
		Format string      `json:"format"`
		// 是否为用户设置的帖子封面
		IsUserSetCover   bool `json:"is_user_set_cover"`
		common.SizeByte  `json:",inline"`
		common.Rectangle `json:",inline"`
	}
)

type (
	TopicDetail struct {
		ID            meta.Topic    `json:"id"`
		Name          string        `json:"name"`
		CoverURL      string        `json:"cover"`
		ContentType   ContentType   `json:"content_type"`
		GID           rest.GameType `json:"game_id"`
		IsGood        bool          `json:"is_good"`
		IsInteractive bool          `json:"is_interactive"`
		IsTop         bool          `json:"is_top"`
	}
)

type (
	PostContent struct {
		PostID                 runtime.Int           `json:"post_id"`
		AuthorID               runtime.Int           `json:"uid"`
		ReviewID               runtime.Int           `json:"review_id"`
		Subject                string                `json:"subject"`
		BodyHTML               string                `json:"content"`
		BodyStructured         string                `json:"structured_content"`
		CoverURL               string                `json:"cover"`
		Forum                  meta.Forum            `json:"f_forum_id"`
		Topics                 []meta.Topic          `json:"topic_ids"`
		ImageURLs              []string              `json:"images"`
		GID                    rest.GameType         `json:"game_id"`
		MaxFloor               int                   `json:"max_floor"`
		Status                 PostStatus            `json:"post_status"`
		RepublishAuthorization RepublishAuthType     `json:"republish_authorization"`
		ViewStatus             ViewStatus            `json:"view_status"`
		ViewType               ViewType              `json:"view_type"`
		IsProfit               bool                  `json:"is_profit"`
		IsInProfit             bool                  `json:"is_in_profit"`
		IsInteractive          bool                  `json:"is_interactive"`
		IsOriginal             OriginalType          `json:"is_original"`
		IsDeleted              DeleteStatus          `json:"is_deleted"`
		CreatedAt              runtime.UnixTimestamp `json:"created_at"`
	}

	PostStatus struct {
		IsTop      bool `json:"is_top"`
		IsGood     bool `json:"is_good"`
		IsOfficial bool `json:"is_official"`
	}
)

type (
	DeleteStatus uint8
	OriginalType uint8

	// 帖子类型
	ContentType uint8

	ViewStatus uint8
	ViewType   uint8

	// 转载授权类型
	RepublishAuthType uint8
)

const (
	ArticlePic ContentType = 3
	Pic        ContentType = 2
)

type (
	// uploaded video
	VodObject struct {
		CoverURL             string               `json:"cover"`
		Length               VodDuration          `json:"duration"`
		AvailableResolutions []VodResolution      `json:"resolutions"`
		ID                   VodID                `json:"id"`
		ReviewStatus         VodReviewStatus      `json:"review_status"`
		TranscodingStatus    VodTranscodingStatus `json:"transcoding_status"`
		Views                int                  `json:"view_num"`
	}

	// in Millisecond
	VodDuration int64

	VodID string

	VodReviewStatus int

	VodTranscodingStatus int

	VodResolution struct {
		// bit/s
		Bitrate          int    `json:"bitrate"`
		Definition       string `json:"definition"`
		Label            string `json:"label"`
		Format           string `json:"format"`
		URL              string `json:"url"`
		common.SizeByte  `json:",inline"`
		common.Rectangle `json:",inline"`
	}
)
