package handler

import (
	"github.com/povsister/mys-mirai/bot"
	"github.com/povsister/mys-mirai/event/router"
	"github.com/povsister/mys-mirai/mys/api/request/meta"
	"github.com/povsister/mys-mirai/pkg/qqmsg"
	"time"
)

func init() {
	router.Router.RegisterHandler(router.Message, "/删帖/:reason", handlePostDelete)
	router.Router.RegisterHandler(router.Message, "/删帖", handlePostDelete)
	router.Router.RegisterHandler(router.Message, "/移水", handleWaterPostMove)
	router.Router.RegisterHandler(router.Message, "/招募", handleFriendRecruit)
	router.Router.RegisterHandler(router.Message, "/删帖/:reason/禁言/:duration", handleDeleteAndMute)
}

const (
	errNoPidFound = `没有发现米游社帖子链接喵~
提示: 可以回复他人发送的帖子链接
或者将帖子链接一并发出来喵`
)

func handlePostDelete(b *bot.Bot, m interface{}, ps router.Params) {
	mys := b.MByMsg(m)
	if mys == nil {
		return
	}
	gid, pid := qqmsg.ExtractMysPostID(m)
	if pid == 0 {
		b.ReplyStr(m, errNoPidFound)
		return
	}

	l := moderatorLogger("delete", pid, gid, mys.Qid())
	var reason meta.DeleteReason
	if ps != nil {
		r := ps.ByName("reason")
		l.Str("reason", r)
		reason = meta.GuessDeleteReason(r)
	}

	l.Msg("deleting mys post")
	err := mys.Moderator().Post(gid).Delete(pid, meta.DeletePostOptions{Reason: reason})
	takeCareMysError(b, m, err)
}

func handleWaterPostMove(b *bot.Bot, m interface{}, _ router.Params) {
	movePost(b, m, meta.NoForum, meta.YsWaterSlime)
}

func movePost(b *bot.Bot, m interface{}, f meta.Forum, topic ...meta.Topic) {
	mys := b.MByMsg(m)
	if mys == nil {
		return
	}
	gid, pid := qqmsg.ExtractMysPostID(m)
	if pid == 0 {
		b.ReplyStr(m, errNoPidFound)
		return
	}
	moderatorLogger("move", pid, gid, mys.Qid()).
		Str("to_forum", f.Name()).
		Interface("topics", topic).
		Msg("moving mys post")
	err := mys.Moderator().Post(gid).Move(pid, meta.MovePostOptions{To: f, WithTopics: topic})
	takeCareMysError(b, m, err)
}

func handleDeleteAndMute(b *bot.Bot, m interface{}, ps router.Params) {
	mys := b.MByMsg(m)
	if mys == nil {
		return
	}
	gid, pid := qqmsg.ExtractMysPostID(m)
	if pid == 0 {
		b.ReplyStr(m, errNoPidFound)
		return
	}
	reason := meta.GuessDeleteReason(ps.ByName("reason"))
	dur, err := time.ParseDuration(ps.ByName("duration"))
	if err != nil {
		b.ReplyStr(m, runtimeErr(err))
		return
	}
	pd, err := mys.Post().Info(gid).Get(pid, meta.GetPostOptions{})
	if err != nil {
		b.ReplyStr(m, runtimeErr(err))
		return
	}

	uid := pd.Post.User.UID.Int()
	moderatorLogger("silent", pid, gid, mys.Qid()).
		Int("uid", uid).
		Str("nickname", pd.Post.User.Nickname).
		Msg("silenting mys user")

	err = mys.Moderator().User(gid).Silence(uid, meta.SilenceUserOptions{Global: dur})
	if err != nil {
		b.ReplyStr(m, "禁言好像出了点问题喵~\n"+err.Error())
		return
	}
	moderatorLogger("delete", pid, gid, mys.Qid()).
		Int("uid", uid).
		Str("nickname", pd.Post.User.Nickname).
		Str("reason", ps.ByName("reason")).
		Msg("deleting mys post")
	err = mys.Moderator().Post(gid).Delete(pid, meta.DeletePostOptions{Reason: reason})
	takeCareMysError(b, m, err)
}

func handleFriendRecruit(b *bot.Bot, m interface{}, _ router.Params) {
	movePost(b, m, meta.NoForum, meta.YsFriendRecruit)
}
