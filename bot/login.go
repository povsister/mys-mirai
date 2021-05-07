package bot

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/povsister/mys-mirai/pkg/log"
	"github.com/povsister/mys-mirai/pkg/util/fs"
	"github.com/povsister/mys-mirai/pkg/util/term"
)

var (
	ErrSessionNotExist = errors.New("session file not exists")
	ErrReadSessionFile = errors.New("failed to read session file")

	ErrLoginFailed = errors.New("login failed")
)

func (b *Bot) Login() (err error) {
	b.lm.mu.Lock()
	defer b.lm.mu.Unlock()
	if b.c.Online {
		return
	}
	if err = b.lm.login(); err == nil {
		b.lm.alreadyLogin = true
		b.lm.setupReconnect()
		b.registerKnownEvents()
		b.c.SetOnlineStatus(client.StatusOnline)
	}
	return
}

// loginManager 负责登录以及断线重连相关的活
type loginManger struct {
	bot *Bot
	// 是否已经登录 用来区分重连
	alreadyLogin bool
	// 重连次数限制
	reConnLimit int
	// 保证同时只有一次登录或者断线重连
	mu   sync.Mutex
	once sync.Once
}

func (l *loginManger) setupReconnect() {
	l.once.Do(func() {
		l.bot.c.OnDisconnected(func(c *client.QQClient, e *client.ClientDisconnectedEvent) {
			l.mu.Lock()
			defer l.mu.Unlock()

			if c.Online {
				return
			}
			log.Warn().Msgf("Bot[%d] 已离线: %s", l.bot.c.Uin, e.Message)
			times := 1
			for times <= l.reConnLimit {
				log.Warn().Msgf("将在 %d 秒后自动重连", 10*times)
				time.Sleep(10 * time.Duration(times) * time.Second)
				log.Warn().Msg("尝试重连 ...")
				if err := l.login(); err == nil {
					return
				}
				times++
			}
		})
	})
}

// must hold mu before calling
func (l *loginManger) login() (err error) {
	// save login token
	defer func() {
		if err != nil {
			return
		}
		log.Info().Msgf("Bot[%d] 登录成功", l.bot.c.Uin)
		if err2 := l.bot.saveSessionFile(l.bot.c.GenToken()); err2 != nil {
			log.Warn().Err(err2).Msg("无法保存 session.token")
			return
		}
	}()

	if err = l.sessionResume(); err == nil {
		return
	}
	// 非文件系统错误 log error
	if err != ErrSessionNotExist && !errors.Is(err, ErrReadSessionFile) {
		log.Warn().Err(err).Msg("无法使用 session.token 恢复登录 尝试正常登录流程")
	}

	if err = l.commonLogin(); err == nil {
		return
	}

	return
}

// 尝试使用 session token 恢复登录.
func (l *loginManger) sessionResume() (err error) {
	sf := l.bot.sessionFile()
	if !fs.FileExists(sf) {
		return ErrSessionNotExist
	}
	defer func() {
		if err != nil {
			fs.RemoveFile(sf)
		}
	}()

	token, err := fs.ReadFile(sf)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrReadSessionFile, err)
	}
	return l.bot.c.TokenLogin(token)
}

func (l *loginManger) commonLogin() error {
	resp, err := l.bot.c.Login()
	if err != nil {
		return err
	}
	return l.handleLoginResponse(resp)
}

func (l *loginManger) handleLoginResponse(r *client.LoginResponse) (err error) {
	var input string
	for {
		if err != nil {
			return
		}
		if r.Success {
			return
		}
		if len(r.ErrorMessage) > 0 {
			log.Debug().Str("from", "login").Msgf("err: %s", r.ErrorMessage)
		}
		// 自动重连的时候一般无人值守，所以碰到错误直接退出
		if l.alreadyLogin {
			log.Error().Str("errMsg", r.ErrorMessage).Msg("重连时出现错误 无法自动处理")
			os.Exit(1)
		}
		switch r.Error {
		case client.SliderNeededError:
			log.Info().Msg("登录需要滑条验证码")
			log.Info().Str("verify_url", r.VerifyUrl).Msg("访问URL并获取Ticket")
			input = term.Readline("Ticket: ")
			r, err = l.bot.c.SubmitTicket(input)
		case client.NeedCaptcha:
			log.Info().Msg("登录需要验证码: captcha.jpg")
			fs.MustWriteFile("captcha.jpg", r.CaptchaImage, 0644)
			input = term.Readline("Captcha: ")
			r, err = l.bot.c.SubmitCaptcha(input, r.CaptchaSign)
		case client.SMSNeededError, client.SMSOrVerifyNeededError:
			log.Info().Msgf("登录需要设备锁 按Enter向手机 %s 发送验证码", r.SMSPhone)
			term.Readline("Press ENTER to continue ...")
			if !l.bot.c.RequestSMS() {
				log.Error().Msg("请求验证码失败")
				return ErrLoginFailed
			}
			input = term.Readline("SMS Code: ")
			r, err = l.bot.c.SubmitSMS(input)
		case client.UnsafeDeviceError:
			log.Info().Msg("登录需要设备锁")
			log.Info().Str("verify_url", r.VerifyUrl).Msg("访问URL验证设备并重启Bot")
			return ErrLoginFailed
		case client.TooManySMSRequestError:
			log.Error().Msg("登录触发太多短信验证请求 请稍后重试")
			return ErrLoginFailed
		case client.OtherLoginError, client.UnknownLoginError:
			log.Error().Msgf("登录失败: %s", r.ErrorMessage)
			return ErrLoginFailed
		}
	}
}
