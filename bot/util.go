package bot

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/povsister/mys-mirai/mys/rest"
	"github.com/povsister/mys-mirai/mys/runtime"
)

func IsPrivateMsg(m interface{}) bool {
	_, ok := m.(*message.PrivateMessage)
	return ok
}

type MessageStringer interface {
	ToString() string
}

func ExtractLastTextElement(m MessageStringer) (ret *message.TextElement) {
	DispatchMessage(m, DispatchTo{
		From: "bot.ExtractLastTextElement",
		Group: func(gm *message.GroupMessage) {
			for i := len(gm.Elements) - 1; i >= 0; i-- {
				if gm.Elements[i].Type() != message.Text {
					continue
				}
				ret = gm.Elements[i].(*message.TextElement)
				break
			}
		},
		Private: func(pm *message.PrivateMessage) {
			for i := len(pm.Elements) - 1; i >= 0; i-- {
				if pm.Elements[i].Type() != message.Text {
					continue
				}
				ret = pm.Elements[i].(*message.TextElement)
				break
			}
		},
	})
	return
}

func ExtractMsgSenderUin(m interface{}) (uin int64) {
	DispatchMessage(m, DispatchTo{
		From: "bot.ExtractMsgSenderUin",
		Group: func(gm *message.GroupMessage) {
			uin = gm.Sender.Uin
		},
		Private: func(pm *message.PrivateMessage) {
			uin = pm.Sender.Uin
		},
	})
	return
}

func toString(e message.IMessageElement, force bool) string {
	t := e.Type()
	if t != message.Text && t != message.Reply && t != message.LightApp && !force {
		return ""
	}
	switch em := e.(type) {
	case *message.TextElement:
		return em.Content
	case *message.ReplyElement:
		strs := make([]string, 0, len(em.Elements))
		for _, im := range em.Elements {
			strs = append(strs, toString(im, force))
		}
		return strings.Join(strs, " ")
	case *message.LightAppElement:
		return em.Content
	default:
		if s, ok := e.(MessageStringer); ok {
			return s.ToString()
		}
	}
	return ""
}

var postIdRgx = regexp.MustCompile(".*bbs\\.mihoyo\\.com[\\\\]{0,2}/([a-z0-9]{2,3}).*article[\\\\]{0,2}/([1-9][0-9]*).*")

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

type LightAppShare struct {
	App    string `json:"app"`
	Config struct {
		AutoSize bool                  `json:"autosize"`
		CTime    runtime.UnixTimestamp `json:"ctime"`
		Forward  bool                  `json:"forward"`
		Token    string                `json:"token"`
		Type     string                `json:"type"`
	} `json:"config"`
	Descr string `json:"desc"`
	Extra struct {
		AppType int   `json:"app_type"`
		AppID   int   `json:"appid"`
		MsgSeq  int64 `json:"msg_seq"`
		Uin     int64 `json:"uin"`
	} `json:"extra"`
	Meta struct {
		View map[string]struct {
			Action         string `json:"action"`
			AndroidPkgName string `json:"android_pkg_name"`
			AppType        int    `json:"app_type"`
			AppID          int64  `json:"appid"`
			Descr          string `json:"desc"`
			JumpURL        string `json:"jump_url"`
			Preview        string `json:"preview"`
			SourceIcon     string `json:"source_icon"`
			SourceURL      string `json:"source_url"`
			Tag            string `json:"tag"`
			Title          string `json:"title"`
		} `json:",inline"`
	} `json:"meta"`
	Prompt string `json:"prompt"`
	Ver    string `json:"ver"`
	View   string `json:"view"`
}

// 尽量把reply elements的内容替换成全的
func enrichGroupMsgFromCache(b *Bot, m *message.GroupMessage) {
	for _, e := range m.Elements {
		if replyE, ok := e.(*message.ReplyElement); ok {
			if prev := b.LookupGroupMessage(m.GroupCode, replyE.ReplySeq); prev != nil {
				replyE.Elements = prev.Elements
			}
		}
	}
}

func enrichPrivateMsgFromCache(b *Bot, m *message.PrivateMessage) {
	for _, e := range m.Elements {
		if replyE, ok := e.(*message.ReplyElement); ok {
			if prev := b.LookupPrivateMessage(m.Sender.Uin, replyE.ReplySeq); prev != nil {
				replyE.Elements = prev.Elements
			}
		}
	}
}

func (b *Bot) ExtractMysPostID(m interface{}) (rest.GameType, int) {
	switch msg := m.(type) {
	case *message.GroupMessage:
		enrichGroupMsgFromCache(b, msg)
		return extractMysPostIdFromElems(msg.Elements)
	case *message.PrivateMessage:
		enrichPrivateMsgFromCache(b, msg)
		return extractMysPostIdFromElems(msg.Elements)
	default:
		logger.Error().Msgf("can not extract uin. unknown message type %T", m)
	}
	return rest.NoGame, 0
}

type DispatchTo struct {
	// debug message.
	// indicate where the message from
	From string
	// GroupMessage callback
	Group func(gm *message.GroupMessage)
	// PrivateMessage callback
	Private func(pm *message.PrivateMessage)
	// Unknown type message callback
	Default func(m interface{})
}

func DispatchMessage(m interface{}, o DispatchTo) {
	switch typedMsg := m.(type) {
	case *message.GroupMessage:
		if o.Group != nil {
			o.Group(typedMsg)
		}
	case *message.PrivateMessage:
		if o.Private != nil {
			o.Private(typedMsg)
		}
	default:
		if o.Default != nil {
			o.Default(m)
		} else {
			l := logger.Error()
			if len(o.From) > 0 {
				l.Str("requested_by", o.From)
			}
			l.Msgf("can not dispatch message. unknown message type %T", m)
		}
	}
}
