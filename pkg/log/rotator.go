package log

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

var (
	ErrLogFileClosed = errors.New("rotator: log file has been closed")
)

const (
	// default rotate Frequency
	Daily RotateFrequency = iota
	Hourly
)

type RotateFrequency uint8

type Rotator struct {
	// the timer to trigger next rotation
	t *time.Timer
	f RotateFrequency
	// max copies of rotated files
	max int

	// the file path
	path string
	// guards the following while rotating
	mu sync.RWMutex
	// the log file fd
	fd     io.WriteCloser
	closed bool
}

func NewRotator(path string) (*Rotator, error) {
	r := &Rotator{path: path, max: 30}
	if err := r.openFile(); err != nil {
		return nil, err
	}
	r.setupRotateTimer()
	return r, nil
}

func (r *Rotator) SetMaxCopies(m int) {
	r.max = m
}

func (r *Rotator) openFile() (err error) {
	r.fd, err = os.OpenFile(r.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	return
}

func (r *Rotator) SetRotateFrequency(f RotateFrequency) {
	r.f = f
	r.setupRotateTimer()
}

func (r *Rotator) setupRotateTimer() {
	if r.t != nil {
		r.t.Stop()
	}
	now := time.Now()
	var dur time.Duration

	switch r.f {
	case Hourly:
		// 下一个小时整
		dur = now.Truncate(time.Hour).Add(time.Hour).Sub(now)
	case Daily:
		nextDay := now.Truncate(time.Hour).Add(24 * time.Hour)
		// 下一天的0点
		dur = nextDay.Add(-time.Duration(nextDay.Hour()) * time.Hour).Sub(now)
	default:
		panic("rotator: unsupported RotateFrequency")
	}
	r.t = time.AfterFunc(dur, r.rotate)
}

func (r *Rotator) rotate() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return
	}
	_ = r.fd.Close()
	err := os.Rename(r.path, r.rotatedFileName())
	if err != nil {
		stdErr("rotator: failed to rotate file %q %v\n", r.path, err)
		r.closed = true
		return
	}
	err = r.openFile()
	if err != nil {
		stdErr("rotator: failed to reopen file %q %v\n", r.path, err)
		r.closed = true
		return
	}
	r.deleteOlderCopies()
	r.setupRotateTimer()
}

func (r *Rotator) deleteOlderCopies() {
	filename := filepath.Base(r.path)
	dir := filepath.Dir(r.path)
	rgx, err := regexp.Compile(filename + "-.+")
	if err != nil {
		stdErr("rotator: failed to clean old log file %v", err)
		return
	}

	fd, err := os.Open(dir)
	if err != nil {
		stdErr("rotator: failed to clean old log file %v", err)
		return
	}
	list, err := fd.ReadDir(-1)
	if err != nil {
		stdErr("rotator: failed to clean old log file %v", err)
		return
	}

	count := 0
	first := ""
	for _, f := range list {
		if !f.IsDir() {
			if f.Name() == filename {
				continue
			}
			if rgx.MatchString(f.Name()) {
				count++
				if len(first) == 0 {
					first = f.Name()
				}
			}
		}
	}

	// remove the first file in directory order
	if count > r.max {
		_ = os.Remove(filepath.Join(dir, first))
	}
}

func (r *Rotator) rotatedFileName() string {
	t := time.Now().Add(-time.Minute)
	var suffix string
	switch r.f {
	case Hourly:
		suffix = fmt.Sprintf("%d-%d-%d_%d", t.Year(), t.Month(), t.Day(), t.Hour())
	default:
		suffix = fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())
	}
	return r.path + "-" + suffix
}

func stdErr(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
}

func (r *Rotator) Write(p []byte) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.closed {
		return 0, ErrLogFileClosed
	}
	return r.fd.Write(p)
}

func (r *Rotator) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.closed = true
	if r.t != nil {
		r.t.Stop()
	}
	return nil
}
