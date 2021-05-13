package meta

import (
	"regexp"

	"github.com/povsister/mys-mirai/mys/rest"
)

type DeleteReason uint8

const (
	// 默认  无原因 向后兼容
	NoReason DeleteReason = 0
	// 包含引战/人身攻击/挂人/攻击他人作品行为
	AttackOthers DeleteReason = 1
	// 复制粘贴跟风刷屏行为
	SPAMorFlood DeleteReason = 2
	// 内容涉嫌侵权
	UnauthorizedRepublish DeleteReason = 3
	// 包含交易、账号转送、广告内容
	AccountTrade DeleteReason = 4
	// 包含其他游戏、社区平台、App相关内容及链接等引流行为
	OuterLink DeleteReason = 5
	// 发布与版区对应游戏无关内容
	UnrelatedContent DeleteReason = 6
	// 包含未证实/虚假信息
	FakeNewsOrReveals DeleteReason = 7
	// 包含违反游戏用户协议内容
	GameRuleViolation DeleteReason = 8
	// 违反社区用户服务协议
	ForumRuleViolation DeleteReason = 9
	// 包含违反国家法律法规政策的内容
	LawViolation DeleteReason = 10
	// 包含违反公序良俗的内容
	MoralityViolation DeleteReason = 11
)

type DeletePostOptions struct {
	Reason DeleteReason
}

func (dpo DeletePostOptions) Apply(r *rest.Request) *rest.Request {
	if dpo.Reason != NoReason {
		r.BodyKV("config_id", dpo.Reason)
	}
	return r
}

var guessRgx = map[DeleteReason]*regexp.Regexp{
	AttackOthers:          regexp.MustCompile(".*引战|挂人|攻击.*"),
	SPAMorFlood:           regexp.MustCompile(".*跟风|刷屏|复制粘贴.*"),
	UnauthorizedRepublish: regexp.MustCompile(".*侵权|未授权.*"),
	AccountTrade:          regexp.MustCompile(".*交易|送号|广告.*"),
	OuterLink:             regexp.MustCompile(".*引流|外部链接.*"),
	UnrelatedContent:      regexp.MustCompile(".*无关|错区.*"),
	FakeNewsOrReveals:     regexp.MustCompile(".*内鬼|未证实|未实.*"),
	GameRuleViolation:     regexp.MustCompile(".*游戏用户协议.*"),
	ForumRuleViolation:    regexp.MustCompile(".*色图|社区用户协议|片哥.*"),
	LawViolation:          regexp.MustCompile(".*违法|国家法律.*"),
	MoralityViolation:     regexp.MustCompile(".*道德|公序良俗.*"),
}

func GuessDeleteReason(s string) (ret DeleteReason) {
	for r, rgx := range guessRgx {
		if rgx.MatchString(s) {
			ret = r
			break
		}
	}
	return
}

type MovePostOptions struct {
	To         Forum
	WithTopics []Topic
}

func (mpo MovePostOptions) Apply(r *rest.Request) *rest.Request {
	r.BodyKV("f_forum_id", mpo.To)
	if len(mpo.WithTopics) > 0 {
		r.BodyKV("topic_ids", mpo.WithTopics)
	}
	return r
}

type RemoveTopicOptions struct {
	Remove []Topic
}

func (rto RemoveTopicOptions) Apply(r *rest.Request) *rest.Request {
	r.BodyKV("topic_ids", rto.Remove)
	return r
}

type GetPostOptions struct {
}

func (gpo GetPostOptions) Apply(r *rest.Request) *rest.Request {
	return r
}
