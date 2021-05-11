package rest

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

type Interface interface {
	Verb(verb string) *Request
	Post() *Request
	Get() *Request
	Put() *Request
	Delete() *Request
	Head() *Request

	Cookies() []*http.Cookie
}

type RESTClient struct {
	counter uint64
	base    url.URL
	c       *resty.Client
}

type Config struct {
	Cookie []*http.Cookie
}

func NewConfig(ck []*http.Cookie) *Config {
	return &Config{Cookie: ck}
}

func NewRESTClient(base url.URL, cfg *Config) *RESTClient {
	c := resty.New()
	c.SetRetryCount(1)
	c.SetTimeout(5 * time.Second)

	if len(cfg.Cookie) > 0 {
		c.SetCookies(cfg.Cookie)
	}
	if !strings.HasSuffix(base.Path, "/") {
		base.Path += "/"
	}

	restC := &RESTClient{
		base: base,
		c:    c,
	}
	restC.setupMiddleware()
	restC.setupCommonHeaders()

	return restC
}

func (c *RESTClient) setupCommonHeaders() {
	uuidGen := uuid.NewString()
	c.c.SetHeaders(map[string]string{
		"Accept":            "application/json",
		"Accept-Encoding":   "gzip, deflate, br",
		"Accept-Language":   "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7",
		"DNT":               "1",
		"Origin":            "https://m.bbs.mihoyo.com",
		"Referer":           "https://m.bbs.mihoyo.com/",
		"User-Agent":        "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
		"x-rpc-app_version": MysAppVersion,
		"x-rpc-client_type": MysClientType,
		"x-rpc-device_id":   uuidGen,
	})
	c.c.SetCookie(&http.Cookie{Name: "_MHYUUID", Value: uuidGen, Domain: ".mihoyo.com"})
}

func (c *RESTClient) setupMiddleware() {
	c.c.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		dsGenerator(request)
		return nil
	})
}

const (
	MysAppVersion = "2.7.0"
	DsSalt        = "14bmu1mz0yuljprsfgpvjh3ju2ni468r"
	MysClientType = "5"

	alphabet = "abcdefghijkmlnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var gRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func randRev(n int) []byte {
	ret := make([]byte, 0, n)
	for i := 1; i <= n; i++ {
		ret = append(ret, alphabet[gRand.Intn(len(alphabet))])
	}
	return ret
}

func dsGenerator(req *resty.Request) {
	now := strconv.FormatInt(time.Now().Unix(), 10)
	rev := randRev(6)
	vals := make([][]byte, 3)
	vals = append(vals, []byte("salt="+DsSalt))
	vals = append(vals, []byte("t="+now))
	vals = append(vals, append([]byte("r="), rev...))

	hash := md5.Sum(bytes.Join(vals, []byte("&")))
	ds := now + "," + string(rev) + "," + fmt.Sprintf("%x", hash)
	req.SetHeader("DS", ds)
}

// return the next request ID. Auto incremental.
func (c *RESTClient) NextRequestID() uint64 {
	return atomic.AddUint64(&c.counter, 1)
}
