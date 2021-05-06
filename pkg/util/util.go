package util

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	globalSignalTerm     chan struct{}
	globalSignalTermOnce = sync.Once{}
)

func SetupSignalHandler() <-chan struct{} {
	globalSignalTermOnce.Do(func() {
		globalSignalTerm = make(chan struct{})
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		close(globalSignalTerm)
		<-c
		os.Exit(1)
	})
	return globalSignalTerm
}
