package rest

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/povsister/mys-mirai/mys/runtime"
	"github.com/povsister/mys-mirai/pkg/log"
)

var logger log.Logger

func init() {
	logger = log.SubLogger("mys.restClient")
}

type bodyBuilder map[string]interface{}

func newBodyBuilder() bodyBuilder {
	return make(map[string]interface{})
}

func (b bodyBuilder) Set(k string, v interface{}) {
	b[k] = v
}

func (b bodyBuilder) Delete(k string) {
	delete(b, k)
}

type Request struct {
	c *RESTClient
	r *resty.Request

	verb    string
	path    string
	params  url.Values
	headers map[string]string
	body    bodyBuilder

	gid GameType

	// error detained
	err error
}

type RequestOptions interface {
	Apply(r *Request) *Request
}

func (r *Request) Use(o RequestOptions) *Request {
	return o.Apply(r)
}

func (r *Request) GID(id GameType) *Request {
	r.gid = id
	return r
}

func (r *Request) Path(path string) *Request {
	r.path = path
	return r
}

func (r *Request) ParamAdd(k, v string) *Request {
	if r.params == nil {
		r.params = url.Values{}
	}
	r.params.Add(k, v)
	return r
}

func (r *Request) ParamSet(k, v string) *Request {
	if r.params == nil {
		r.params = url.Values{}
	}
	r.params.Set(k, v)
	return r
}

// set request header
func (r *Request) Header(k, v string) *Request {
	if r.headers == nil {
		r.headers = make(map[string]string)
	}
	r.headers[k] = v
	return r
}

// k-v request body. Auto marshall to json
func (r *Request) BodyKV(k string, v interface{}) *Request {
	if r.body == nil {
		r.body = newBodyBuilder()
	}
	r.body.Set(k, v)
	return r
}

// overrides the BodyKV
func (r *Request) Body(obj interface{}) *Request {
	r.r.SetBody(obj)
	return r
}

func (r *Request) Do() *Result {
	if r.err != nil {
		return &Result{err: r.err}
	}
	thisID := r.c.NextRequestID()
	// request logger
	l := logger.Trace().Str("method", r.verb)

	r.r.Method = r.verb
	// set request body
	if r.body != nil {
		if r.gid != NoGame {
			r.body.Set("gids", r.gid)
		}
		r.r.SetBody(r.body)
		l.Interface("body", r.body)
	}
	// set request param
	if len(r.params) > 0 && r.gid != NoGame {
		r.params.Set("gids", string(r.gid))
		r.r.QueryParam = r.params
		l.Str("query", r.r.QueryParam.Encode())
	}
	// compose request url
	u := url.URL{
		Scheme: r.c.base.Scheme,
		Host:   r.c.base.Host,
		Path:   r.c.base.Path + strings.TrimLeft(r.path, "/"),
	}
	r.r.URL = u.String()
	l.Str("url", r.r.URL)
	// set additional headers if any
	if r.headers != nil {
		r.r.SetHeaders(r.headers)
		l.Interface("header", r.headers)
	}
	// send out the request and logger
	l.Uint64("request_id", thisID).Msg("request sent")
	resp, err := r.r.Send()
	if err != nil {
		logger.Trace().Uint64("request_id", thisID).Err(err).
			Msg("request sending error")
		return &Result{resp: resp, err: err}
	}
	return r.transformResponse(thisID, resp)
}

func (r *Request) transformResponse(rid uint64, resp *resty.Response) *Result {
	decoder := runtime.NegotiateDecoder(resp.Header().Get("Content-Type"))
	body := resp.Body()
	logger.Trace().Uint64("request_id", rid).
		Bytes("body", body).Int("status", resp.StatusCode()).
		Msgf("%s response received", resp.Status())
	return &Result{resp: resp, body: resp.Body(), decoder: decoder}
}

type Result struct {
	resp *resty.Response
	body []byte
	err  error

	decoder runtime.Decoder
}

// 只返回 runtime error. 不检查 application error.
// 例如 对于 retcode != 0 的情况不做检查
func (r *Result) Into(obj runtime.Object) error {
	if r.err != nil {
		return r.err
	}
	if r.decoder == nil {
		return fmt.Errorf("no decoder specified for content-type %q", r.resp.Header().Get("Content-Type"))
	}
	if len(r.body) == 0 {
		return fmt.Errorf("empty response with status %d", r.resp.StatusCode())
	}

	err := r.decoder.Decode(r.body, obj)
	if err != nil {
		logger.Warn().Err(err).Bytes("body", r.body).Msg("unable to decode response body")
		return err
	}
	return nil
}

// 优先返回 runtime error.
// 其次检查 application level error.
// 例如 对于 retcode != 0 的响应, 返回 retcode: message 作为错误信息
func (r *Result) Error() error {
	if len(r.body) == 0 || r.decoder == nil {
		return r.err
	}

	// trying to decode the error
	obj := &runtime.ObjectMeta{}
	err := r.decoder.Decode(r.body, obj)
	if err != nil {
		logger.Warn().Err(err).Bytes("body", r.body).Msg("unable to decode response body")
		return r.err
	}
	if obj.Retcode() != runtime.OK {
		return obj
	}
	return nil
}

func IsApplicationErr(err error) bool {
	_, ok := err.(runtime.Object)
	return ok
}
