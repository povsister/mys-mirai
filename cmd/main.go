package main

import (
	"github.com/povsister/mys-mirai/pkg/log"
)

func main() {
	log.InitLogger()

	log.Info().Msg("log is good")
}
