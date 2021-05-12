package post

import (
	"github.com/povsister/mys-mirai/mys/api/response/forum"
	"github.com/povsister/mys-mirai/mys/api/response/user"
	"github.com/povsister/mys-mirai/mys/runtime"
)

type (
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
		CoverURL  string               `json:"cover"`
		Forum     forum.ForumInfoBasic `json:"forum"`
		User      user.UserInfoBasic   `json:"user"`
		Stat      PostStat             `json:"stat"`
		VodList   []VodObject          `json:"vod_list"`
		VoteCount int                  `json:"vote_count"`
	}

	PostStat struct {
		Bookmarks int `json:"bookmark_num"`
		Likes     int `json:"like_num"`
		Replies   int `json:"reply_num"`
		Views     int `json:"view_num"`
	}
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

	VodReviewStatus int8

	VodTranscodingStatus int8

	// in Bytes
	VodSize string

	VodResolution struct {
		// bit/s
		Bitrate    int     `json:"bitrate"`
		Definition string  `json:"definition"`
		Label      string  `json:"label"`
		Format     string  `json:"format"`
		Width      uint16  `json:"width"`
		Height     uint16  `json:"height"`
		URL        string  `json:"url"`
		Size       VodSize `json:"size"`
	}
)
