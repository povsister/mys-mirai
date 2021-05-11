package mys

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"

	"github.com/povsister/mys-mirai/mys/api/request/meta"
	"github.com/povsister/mys-mirai/mys/rest"
	"github.com/povsister/mys-mirai/mys/runtime"
	"github.com/povsister/mys-mirai/pkg/log"
	"github.com/povsister/mys-mirai/pkg/util/fs"
	"github.com/povsister/mys-mirai/resources"
)

var (
	logger         log.Logger
	validCookieRgx = regexp.MustCompile("^([0-9]+)\\.cookie$")

	CookieFolder = "cookies"
)

func init() {
	logger = log.SubLogger("mys.userManager")
}

// 根据 cookie 生成对应 mys client.
// 并与 QQ号 对应关联.
// 包括解析、自动持久化 cookie 等等
type UserManager struct {
	// guards cs
	mu sync.RWMutex
	cs map[int64]*Clientset
}

func NewUserManager() *UserManager {
	m := &UserManager{
		cs: map[int64]*Clientset{},
	}
	err := fs.MkDir(CookieFolder)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to mkdir %q", CookieFolder)
	}
	m.scanEmbedUser()
	m.scanPersistUser()
	return m
}

func (m *UserManager) scanEmbedUser() {
	embedUser, err := resources.FS.ReadDir("cookie")
	if err != nil {
		logger.Warn().Err(err).Msg("failed to list embed cookies")
		return
	}
	for _, u := range embedUser {
		if !u.IsDir() {
			if qid := parseQid(u.Name()); qid != 0 {
				cookie, err := resources.FS.ReadFile("cookie/" + u.Name())
				if err != nil {
					logger.Warn().Err(err).Msg("failed to read embed cookie")
					continue
				}
				m.AddByNetscapeCookie(qid, cookie, false)
			}
		}
	}
}

func parseQid(fName string) int64 {
	if match := validCookieRgx.FindStringSubmatch(fName); len(match) > 0 {
		qqUin, _ := strconv.ParseInt(match[1], 10, 64)
		return qqUin
	}
	return 0
}

func (m *UserManager) AddByNetscapeCookie(qid int64, cookie []byte, persistCookie bool) {
	cks := fs.ParseNetscapeCookie(cookie)
	if cks != nil {
		config := rest.NewConfig(cks)
		config.Qid = qid
		m.mu.Lock()
		defer m.mu.Unlock()
		m.cs[qid] = NewClient(config)
		if persistCookie {
			m.persistCookie(qid, cks)
		}
		logger.Info().Int64("qid", qid).Msg("added mys client from persist cookie")
	}
}

// must be called with mu hold
func (m *UserManager) persistCookie(qid int64, ck []*http.Cookie) {
	file := filepath.Join(CookieFolder, fmt.Sprintf("%d.cookie", qid))
	data, err := fs.ToNetscapeCookie(ck)
	if err != nil {
		logger.Error().Err(err).Int64("qid", qid).Msg("can not marshall cookie to Netscape format")
		return
	}
	err = fs.WriteFile(file, data, 0640)
	if err != nil {
		logger.Error().Err(err).Int64("qid", qid).Msg("can not persist cookie")
		return
	}
	logger.Info().Int64("qid", qid).Str("file", file).Msg("successfully persists cookie")
}

func (m *UserManager) erasePersistCookie(qid int64) {
	file := filepath.Join(CookieFolder, fmt.Sprintf("%d.cookie", qid))
	fs.RemoveFile(file)
}

func (m *UserManager) scanPersistUser() {
	list, err := fs.ReadDir(CookieFolder)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to list %q folder", CookieFolder)
		return
	}
	for _, f := range list {
		if f.IsDir() {
			continue
		}
		qid := parseQid(f.Name())
		if qid == 0 {
			continue
		}
		data, err := fs.ReadFile(filepath.Join(CookieFolder, f.Name()))
		if err != nil {
			logger.Error().Err(err).Msg("can not read persist cookie")
			continue
		}
		m.AddByNetscapeCookie(qid, data, false)
	}
}

func (m *UserManager) Add(qid int64, c *Clientset, persistCookie bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cs[qid] = c
	logger.Info().Int64("qid", qid).Msg("added mys client")
	if persistCookie {
		m.persistCookie(qid, c.RESTClient().Cookies())
	}
}

func (m *UserManager) Del(qid int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.cs, qid)
	logger.Info().Int64("qid", qid).Msg("removed mys client")
	m.erasePersistCookie(qid)
}

func (m *UserManager) Get(qid int64) *Clientset {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.cs[qid]
}

func (m *UserManager) TestCookie(ck []*http.Cookie) (*Clientset, error) {
	config := rest.NewConfig(ck)
	c := NewClient(config)
	myInfo, err := c.User().Info(rest.NoGame).Get(meta.UserMyself, meta.UserInfoGetOptions{})
	if err != nil {
		return nil, err
	}
	if !myInfo.LoggedIn() {
		if myInfo.Code == runtime.NeedLogin {
			return nil, fmt.Errorf("%w: %s", runtime.ErrLoginNeeded, myInfo.Message)
		}
		return nil, myInfo
	}
	return c, nil
}
