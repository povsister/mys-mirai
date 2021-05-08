package runtime

import "fmt"

type Object interface {
	Retcode() Retcode
	RetMessage() string
}

type Retcode int

const (
	OK Retcode = 0
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
