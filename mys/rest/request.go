package rest

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/povsister/mys-mirai/mys/api/request/meta"
	"github.com/povsister/mys-mirai/mys/runtime"
	"github.com/povsister/mys-mirai/pkg/log"
)

var logger = log.With().Str("from", "mys.rest").Logger()

type BodyBuilder map[string]interface{}

func NewBodyBuilder() BodyBuilder {
	return make(map[string]interface{})
}

func (b BodyBuilder) Set(k string, v interface{}) BodyBuilder {
	b[k] = k
	return b
}

func (b BodyBuilder) Delete(k string) BodyBuilder {
	delete(b, k)
	return b
}

type Request struct {
	c *RESTClient
	r *resty.Request

	verb    string
	path    string
	params  url.Values
	headers http.Header
	body    BodyBuilder

	gid meta.GameType

	err error
}

func (r *Request) GID(id meta.GameType) *Request {
	r.gid = id
	return r
}

func (r *Request) Path(path string) *Request {
	r.path = path
	return r
}

// k-v request body. Auto marshall to json
func (r *Request) BodyKV(k string, v interface{}) *Request {
	if r.body == nil {
		r.body = NewBodyBuilder()
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
	r.r.Method = r.verb
	if r.body != nil {
		if r.gid != meta.NoForum {
			r.body.Set("gid", r.gid)
		}
		r.r.SetBody(r.body)
	}
	resp, err := r.r.Send()
	if err != nil {
		return &Result{err: err}
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
		return r.Error()
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
