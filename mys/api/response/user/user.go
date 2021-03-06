package user

import (
	"github.com/povsister/mys-mirai/mys/api/meta"
	"github.com/povsister/mys-mirai/mys/rest"
	"github.com/povsister/mys-mirai/mys/runtime"
)

type (
	// /user/wapi/getUserFullInfo?gids=2
	FullUserInfoResponse struct {
		runtime.ObjectMeta `json:",inline"`
		FullUserInfoData   `json:"data"`
	}

	FullUserInfoData struct {
		Privilege     AllPrivilege   `json:"auth_relations"`
		FanRelation   FanStatus      `json:"follow_relation"`
		IsCreator     bool           `json:"is_creator"`
		HasCollection bool           `json:"is_has_collection"`
		Blacklisted   bool           `json:"is_in_blacklist"`
		UserInfo      UserInfoDetail `json:"user_info"`
	}
)

func (fui *FullUserInfoResponse) LoggedIn() bool {
	return fui.Code == runtime.OK
}

const (
	Unknown Gender = 0
	Male    Gender = 1
	Female  Gender = 2
)

const (
	// 适用于整个 gid 分区
	Generic PrivilegeType = 0
	// 只适用于某个版块
	ForumOnly PrivilegeType = 1
)

const (
	SilenceUser         Permission = "silent"
	ForbidUser          Permission = "forbid"
	ChatWithoutFollowed Permission = "chat_skip_follow_limit"

	// moderators

	// 置顶
	StickPost Permission = "top_post"
	// 加精
	PickPost     Permission = "good_post"
	DeletePost   Permission = "delete_post"
	HidePost     Permission = "hide_post"
	ViewHidePost Permission = "view_hide_post"
	CreatePost   Permission = "release_post"
	MovePost     Permission = "move_post"
	// 版区禁言
	ForumSilence    Permission = "forum_silent"
	DeleteReply     Permission = "delete_reply"
	RemovePostTopic Permission = "remove_in_topic"
)

const (
	Official   CertificationType = 1
	Individual CertificationType = 2
)

const (
	HistoryView MysFeature = "enable_history_view"
	Recommend   MysFeature = "enable_recommend"
)

type FanStatus struct {
	IamFan  bool `json:"is_following"`
	IsMyFan bool `json:"is_followed"`
}

type (
	AllPrivilege []Privilege

	Privilege struct {
		Forum       meta.Forum    `json:"business_id"`
		Type        PrivilegeType `json:"business_type"`
		GID         rest.GameType `json:"game_id"`
		Permissions []Permission  `json:"permissions"`
	}

	PrivilegeType uint8

	Permission string
)

type (
	UserInfoDetail struct {
		UserInfoBasic  `json:",inline"`
		Achieve        Achieve               `json:"achieve"`
		Certifications []CertificationDetail `json:"certifications"`
		CommunityInfo  MysInfo               `json:"community_info"`
		Exps           []ExpStatus           `json:"level_exps"`
	}

	UserInfoBasic struct {
		Avatar          Avatar        `json:"avatar"`
		AvatarURL       string        `json:"avatar_url"`
		Certification   Certification `json:"certification"`
		Gender          Gender        `json:"gender"`
		SelfDescription string        `json:"introduce"`
		CurrentExp      ExpStatus     `json:"level_exp"`
		Nickname        string        `json:"nickname"`
		PendantURL      string        `json:"pendant"`
		UID             runtime.Int   `json:"uid"`
		// the following fields may not exist in UserInfoDetail
		FanStatus `json:",inline"`
	}

	Avatar runtime.Int
	Gender int8

	Achieve struct {
		Following        runtime.Int `json:"follow_cnt"`
		FollowCollection runtime.Int `json:"follow_collection_cnt"`
		Follower         runtime.Int `json:"followed_count"`
		GoodPost         runtime.Int `json:"good_post_num"`
		Likes            runtime.Int `json:"like_num"`
		NewFollower      runtime.Int `json:"new_follower_num"`
		PostCreated      runtime.Int `json:"post_num"`
		ReplyCreated     runtime.Int `json:"replypost_num"`
		TopicCreated     runtime.Int `json:"topic_cnt"`
	}
)

type (
	CertificationType uint8

	Certification struct {
		Type  CertificationType `json:"type"`
		Label string            `json:"label"`
	}

	CertificationDetail struct {
		ID              runtime.Int       `json:"id"`
		CertificationID runtime.Int       `json:"certification_id"`
		Label           string            `json:"label"`
		Type            CertificationType `json:"type"`
	}
)

type (
	MysInfo struct {
		EULA               bool                  `json:"agree_status"`
		BannedUntil        runtime.UnixTimestamp `json:"forbid_end_time"`
		ForumSilence       []SilenceInfo         `json:"forum_silent_info"`
		Initialized        bool                  `json:"has_initialized"`
		UpdatedTime        runtime.UnixTimestamp `json:"info_upd_time"`
		Realnamed          bool                  `json:"is_realname"`
		GlobalSilenceUntil runtime.UnixTimestamp `json:"silent_end_time"`
		FeatureEnabled     map[MysFeature]bool   `json:"user_func_status"`
		NoNotify           NotifyDisabled        `json:"notify_disable"`
		Privacy            PrivacyInvisible      `json:"privacy_invisible"`
	}

	SilenceInfo struct {
		Forum     meta.Forum            `json:"forum_id"`
		ForumName string                `json:"forum_name"`
		Until     runtime.UnixTimestamp `json:"forum_silent_end_time"`
		GID       rest.GameType         `json:"game_id"`
	}

	MysFeature string

	NotifyDisabled struct {
		NewChat      bool `json:"chat"`
		NewFollower  bool `json:"follow"`
		NewReply     bool `json:"reply"`
		NewSystemMsg bool `json:"system"`
		NewLikes     bool `json:"upvote"`
	}

	PrivacyInvisible struct {
		HideCollection bool `json:"collect"`
		HidePostReply  bool `json:"post"`
		NoPicWatermark bool `json:"watermark"`
	}

	ExpStatus struct {
		TotalExp int           `json:"exp"`
		Level    int           `json:"level"`
		GID      rest.GameType `json:"game_id"`
	}
)
