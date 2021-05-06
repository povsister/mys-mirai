package term

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var console = bufio.NewReader(os.Stdin)

func Readline(prompt ...interface{}) string {
	if len(prompt) > 0 {
		_, _ = fmt.Fprint(os.Stdout, prompt...)
	}
	s, _ := console.ReadString('\n')
	return strings.TrimSpace(s)
}
