package rest

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/povsister/mys-mirai/mys/runtime"
	"github.com/povsister/mys-mirai/pkg/log"
)

var logger = log.With().Str("from", "mys.rest").Logger()

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

	r.r.Method = r.verb
	// set request body
	if r.body != nil {
		if r.gid != NoForum {
			r.body.Set("gid", r.gid)
		}
		r.r.SetBody(r.body)
	}
	// set request param
	if len(r.params) > 0 && r.gid != NoForum {
		r.params.Set("gid", strconv.Itoa(int(r.gid)))
		r.r.QueryParam = r.params
	}
	// compose request url
	u := url.URL{
		Scheme: r.c.base.Scheme,
		Host:   r.c.base.Host,
		Path:   r.c.base.Path + strings.TrimLeft(r.path, "/"),
	}
	r.r.URL = u.String()
	// set additional headers if any
	if r.headers != nil {
		r.r.SetHeaders(r.headers)
	}
	// send out the request
	resp, err := r.r.Send()
	if err != nil {
		return &Result{resp: resp, err: err}
	}
	return r.transformResponse(resp)
}

func (r *Request) transformResponse(resp *resty.Response) *Result {
	decoder := runtime.NegotiateDecoder(resp.Header().Get("Content-Type"))
	return &Result{resp: resp, body: resp.Body(), decoder: decoder}
}

type Result struct {
	resp *resty.Response
	body []byte
	err  error

	decoder runtime.Decoder
}

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
