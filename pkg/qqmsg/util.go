package qqmsg

import (
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/povsister/mys-mirai/mys/rest"
	"github.com/povsister/mys-mirai/pkg/log"
	"regexp"
	"strconv"
	"strings"
)

var logger log.Logger

func init() {
	logger = log.SubLogger("util.qqmsg")
}

func IsPrivateMsg(m interface{}) bool {
	_, ok := m.(*message.PrivateMessage)
	return ok
}

type MessageStringer interface {
	ToString() string
}

func ExtractLastTextElement(m MessageStringer) *message.TextElement {
	switch msg := m.(type) {
	case *message.GroupMessage:
		for i := len(msg.Elements) - 1; i >= 0; i-- {
			if msg.Elements[i].Type() != message.Text {
				continue
			}
			return msg.Elements[i].(*message.TextElement)
		}
	case *message.PrivateMessage:
		for i := len(msg.Elements) - 1; i >= 0; i-- {
			if msg.Elements[i].Type() != message.Text {
				continue
			}
			return msg.Elements[i].(*message.TextElement)
		}
	default:
		logger.Error().Str("msg", m.ToString()).
			Msgf("can not handle message. unknown message type %T", m)
	}
	return nil
}

func ExtractMsgSenderUin(m interface{}) int64 {
	switch msg := m.(type) {
	case *message.GroupMessage:
		return msg.Sender.Uin
	case *message.PrivateMessage:
		return msg.Sender.Uin
	default:
		logger.Error().Msgf("can not extract uin. unknown message type %T", m)
	}
	return 0
}

func toString(e message.IMessageElement, force bool) (ret string) {
	t := e.Type()
	if t != message.Text && t != message.Reply && !force {
		return
	}
	switch em := e.(type) {
	case *message.TextElement:
		ret = em.Content
	case *message.ReplyElement:
		strs := make([]string, len(em.Elements))
		for _, im := range em.Elements {
			strs = append(strs, toString(im, force))
		}
		ret = strings.Join(strs, " ")
	default:
		if s, ok := e.(MessageStringer); ok {
			ret = s.ToString()
		}
	}
	return
}

var postIdRgx = regexp.MustCompile(".*bbs\\.mihoyo\\.com/([a-z0-9]{2,3}).*article/([1-9][0-9]*).*")

func extractMysPostIdFromElems(es []message.IMessageElement) (rest.GameType, int) {
	for _, e := range es {
		s := toString(e, false)
		if m := postIdRgx.FindStringSubmatch(s); len(m) > 0 {
			pid, err := strconv.Atoi(m[2])
			if err != nil {
				logger.Error().Err(err).Msg("can not extract mys post id")
				return rest.NoGame, 0
			}
			return rest.FromGidStr(m[1]), pid
		}
	}
	return rest.NoGame, 0
}

func ExtractMysPostID(m interface{}) (rest.GameType, int) {
	switch msg := m.(type) {
	case *message.GroupMessage:
		return extractMysPostIdFromElems(msg.Elements)
	case *message.PrivateMessage:
		return extractMysPostIdFromElems(msg.Elements)
	default:
		logger.Error().Msgf("can not extract uin. unknown message type %T", m)
	}
	return rest.NoGame, 0
}
