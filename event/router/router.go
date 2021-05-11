// Copyright 2013 Julien Schmidt. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.
//
// Modified by povsister
package router

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/povsister/mys-mirai/bot"
	"github.com/povsister/mys-mirai/pkg/log"
	"github.com/povsister/mys-mirai/pkg/qqmsg"
	"io"
	"strings"
	"sync"
)

var (
	Router      *router
	NoopHandler Handle = func(bot *bot.Bot, msg interface{}, ps Params) {}

	logger log.Logger
)

func init() {
	logger = log.SubLogger("router")

	Router = NewRouter()
}

type Handle func(bot *bot.Bot, msg interface{}, ps Params)

type callbackType string

const (
	// 群聊或者私聊消息
	Message callbackType = "message"
	// 各种事件
	Event callbackType = "event"
	// 邀请加群请求 或者
	Request callbackType = "request"
)

type router struct {
	trees map[callbackType]*node

	paramsPool sync.Pool
	maxParams  uint16
}

func NewRouter() *router {
	return &router{}
}

// get a reusable params form pool
func (r *router) getParams() *Params {
	ps, _ := r.paramsPool.Get().(*Params)
	*ps = (*ps)[0:0] // reset slice
	return ps
}

// put it back into pool
func (r *router) putParams(ps *Params) {
	if ps != nil {
		r.paramsPool.Put(ps)
	}
}

func (r *router) RegisterHandler(t callbackType, path string, handler Handle) {
	if r.trees == nil {
		r.trees = make(map[callbackType]*node)
	}
	root := r.trees[t]
	if root == nil {
		root = new(node)
		r.trees[t] = root
	}

	root.addRoute(path, handler)

	if paramsCount := countParams(path); paramsCount > r.maxParams {
		r.maxParams = paramsCount
	}

	if r.paramsPool.New == nil && r.maxParams > 0 {
		r.paramsPool.New = func() interface{} {
			ps := make(Params, 0, r.maxParams)
			return &ps
		}
	}
}

func (r *router) handle(b *bot.Bot, t callbackType, path string, m interface{}) {
	root := r.trees[t]
	if root == nil {
		logger.Debug().Str("path", path).Msgf("no root handler for %q", Message)
		return
	}

	l := logger.Debug().Str("path", path)
	switch knownMsg := m.(type) {
	case qqmsg.MessageStringer:
		l.Str("msg", knownMsg.ToString())
	}

	handle, ps, tsr := root.getValue(path, r.getParams)
	if handle != nil {
		// do the handling
		if ps != nil {
			l.Interface("params", *ps).Msgf("handling %s", t)
			handle(b, m, *ps)
			r.putParams(ps)
		} else {
			l.Msgf("handling %s", t)
			handle(b, m, nil)
		}
	} else {
		if tsr {
			l.Msg("extra trailing slash found. possible wrong syntax")
			b.ReplyStr(m, fmt.Sprintf(errCommandTSR, strings.TrimRight(path, "/")))
		} else {
			l.Msg("no handler found")
			//b.ReplyStr(m, errCommandNotFound)
		}
	}
}

func readCommandToken(msgByte, text []byte) string {
	idx := bytes.Index(text, []byte("/"))
	if idx == -1 {
		//logger.Debug().Bytes("msg", msgByte).Bytes("text", text).Msg("no command pattern")
		return ""
	}
	path, err := bufio.NewReader(bytes.NewReader(text[idx:])).ReadBytes(' ')
	if err != nil && err != io.EOF {
		logger.Error().Bytes("msg", msgByte).Bytes("text", text).Err(err).
			Msg("failed to parse command token")
		return ""
	}
	return string(bytes.Trim(path, " \r\n"))
}

func (r *router) HandleMessage(b *bot.Bot, m qqmsg.MessageStringer) {
	msgByte := []byte(m.ToString())
	if len(msgByte) == 0 {
		return
	}
	text := qqmsg.ExtractLastTextElement(m)
	if text == nil {
		return
	}
	path := readCommandToken(msgByte, []byte(text.Content))
	if len(path) == 0 {
		return
	}
	r.handle(b, Message, path, m)
}
