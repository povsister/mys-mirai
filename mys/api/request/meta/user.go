package meta

import (
	"time"

	"github.com/povsister/mys-mirai/mys/rest"
)

type ForumSilence struct {
	ID       Forum
	Duration time.Duration
}

type SilenceUserOptions struct {
	//  全米游社禁言时间
	Global time.Duration
	// 单版块禁言
	SilenceForum []ForumSilence
}

func (suo SilenceUserOptions) Apply(r *rest.Request) *rest.Request {
	// 目前不知道是干啥用的
	r = r.BodyKV("source", 1)
	if suo.Global > 0 {
		return r.BodyKV("duration", int(suo.Global.Seconds()))
	}
	var fs []map[string]int
	for _, per := range suo.SilenceForum {
		fs = append(fs, map[string]int{
			"forum_id":              int(per.ID),
			"forum_silent_duration": int(per.Duration.Seconds()),
		})
	}
	return r.BodyKV("forum_silence", fs)
}

type UnSilenceUserOptions struct {
	// 全米游社解禁
	Global bool
	// 单版块解禁  目前这个API不支持多板块同时解禁
	UnSilenceForum Forum
}

func (usuo UnSilenceUserOptions) Apply(r *rest.Request) *rest.Request {
	if usuo.Global {
		return r
	}
	return r.BodyKV("forum_id", usuo.UnSilenceForum)
}

// 自己的uid 用于不填的默认行为
const MySelfUser = 0

type UserInfoGetOptions struct {
}

func (uio UserInfoGetOptions) Apply(r *rest.Request) *rest.Request {
	return r
}
