package handler

import (
	"github.com/povsister/mys-mirai/bot"
	"github.com/povsister/mys-mirai/mys/rest"
	"github.com/povsister/mys-mirai/pkg/log"
	"github.com/rs/zerolog"
)

var logger log.Logger

func init() {
	logger = log.SubLogger("router.handler")
}

func takeCareMysError(b *bot.Bot, originM interface{}, err error) {
	if err == nil {
		b.SendOrReplyStrByMsg(originM, "kira～☆")
		return
	}
	if rest.IsApplicationErr(err) {
		b.SendOrReplyStrByMsg(originM, applicationErr(err))
		return
	}
	b.SendOrReplyStrByMsg(originM, runtimeErr(err))
	return
}

func moderatorLogger(opType string, pid int, gid rest.GameType, qid int64) *zerolog.Event {
	l := logger.Debug().Str("op_type", opType)
	if pid != 0 {
		l.Int("post_id", pid)
	}
	return l.Str("gid", gid.Name()).
		Int64("qid", qid)
}
