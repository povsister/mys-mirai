package runtime

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/povsister/mys-mirai/pkg/log"
)

var logger log.Logger

func init() {
	logger = log.SubLogger("mys.runtime")
}

var (
	ErrLoginNeeded = errors.New("需要登录")
)

type Object interface {
	Retcode() Retcode
	RetMessage() string
}

type Retcode int

const (
	OK        Retcode = 0
	NeedLogin Retcode = -100
)

type ObjectMeta struct {
	Code    Retcode `json:"retcode"`
	Message string  `json:"message"`
}

func (rm *ObjectMeta) Retcode() Retcode {
	return rm.Code
}

func (rm *ObjectMeta) RetMessage() string {
	return rm.Message
}

func (rm *ObjectMeta) Error() string {
	return fmt.Sprintf("%d: %s", rm.Code, rm.Message)
}

type Int string

func (n Int) Int() int {
	ret, err := strconv.ParseInt(string(n), 10, 64)
	if err != nil {
		logger.Error().Err(err).Msgf("can not convert %q to int", n)
		return 0
	}
	return int(ret)
}

type UnixTimestamp int64
