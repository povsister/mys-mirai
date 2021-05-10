package bot

import (
	"errors"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/povsister/mys-mirai/pkg/log"
)

type miraiLogLevel string

const (
	Trace miraiLogLevel = "TRACE"
	Debug miraiLogLevel = "DEBUG"
	Info  miraiLogLevel = "INFO"
	Warn  miraiLogLevel = "WARNING"
	Error miraiLogLevel = "ERROR"
)

var ErrUnknownLogLevel = errors.New("unknown mirai loglevel")

// 设置来自mirai core的log转发
func (b *Bot) setupLogHandler() {
	sublogger := log.SubLogger("mirai")
	b.c.OnLog(func(c *client.QQClient, e *client.LogEvent) {
		switch miraiLogLevel(e.Type) {
		case Trace:
			sublogger.Trace().Msg(e.Message)
		case Debug:
			sublogger.Debug().Msg(e.Message)
		case Info:
			sublogger.Info().Msg(e.Message)
		case Warn:
			sublogger.Warn().Msg(e.Message)
		case Error:
			sublogger.Error().Msg(e.Message)
		default:
			sublogger.Warn().Err(ErrUnknownLogLevel).Msg(e.Message)
		}
	})
}
