package runtime

import (
	"errors"
	"fmt"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/povsister/mys-mirai/pkg/log"
)

var (
	logger log.Logger

	NumberCodec = &numberCodec{}
)

func init() {
	logger = log.SubLogger("mys.runtime")
}

var (
	ErrLoginNeeded = errors.New("需要登录")
)

type Retcode int

const (
	OK        Retcode = 0
	NeedLogin Retcode = -100
)

type Object interface {
	Retcode() Retcode
	RetMessage() string
}

type ObjectMeta struct {
	Code    Retcode `json:"retcode"`
	Message string  `json:"message"`
}

func (rm ObjectMeta) Retcode() Retcode {
	return rm.Code
}

func (rm ObjectMeta) RetMessage() string {
	return rm.Message
}

func (rm ObjectMeta) Error() string {
	return fmt.Sprintf("%d: %s", rm.Code, rm.Message)
}

// Number codec for json encode/decode
type numberCodec struct {
}

func (codec *numberCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	switch iter.WhatIsNext() {
	case jsoniter.StringValue:
		*((*jsoniter.Number)(ptr)) = jsoniter.Number(iter.ReadString())
	case jsoniter.NilValue:
		if ok := iter.ReadNil(); !ok {
			iter.ReportError("readNil", "expect null")
		}
		*((*jsoniter.Number)(ptr)) = ""
	default:
		*((*jsoniter.Number)(ptr)) = jsoniter.Number(iter.ReadNumber().String())
	}
}

func (codec *numberCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	number := *((*jsoniter.Number)(ptr))
	if len(number) == 0 {
		stream.WriteInt(0)
	} else {
		stream.WriteRaw(string(number))
	}
}

func (codec *numberCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return len(*((*jsoniter.Number)(ptr))) == 0
}

type (
	// an Int representation
	Int jsoniter.Number
	// unixTimeStamp in second
	UnixTimestamp jsoniter.Number
	// "2006-01-02T15:04:05Z07:00"
	TimeRFC3339 string
)

func init() {
	jsoniter.RegisterTypeDecoder("runtime.Int", NumberCodec)
	jsoniter.RegisterTypeEncoder("runtime.Int", NumberCodec)
}

func (n Int) Int() int {
	ret, err := jsoniter.Number(n).Int64()
	if err != nil {
		logger.Error().Err(err).Msgf("can not convert %q to int", string(n))
		return 0
	}
	return int(ret)
}

func (ut UnixTimestamp) Time() time.Time {
	num, err := jsoniter.Number(ut).Int64()
	if err != nil {
		logger.Error().Err(err).Msgf("can not parse unixtimestamp from %q", string(ut))
		return time.Unix(0, 0)
	}
	return time.Unix(num, 0)
}

func (t TimeRFC3339) Time() time.Time {
	tf, err := time.Parse(time.RFC3339, string(t))
	if err != nil {
		logger.Error().Err(err).Msgf("can not parse RFC3339 time from %q", string(t))
		return time.Unix(0, 0)
	}
	return tf
}
