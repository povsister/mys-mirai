package rest

import "net/http"

func (c *RESTClient) Verb(verb string) *Request {
	req := &Request{c: c, r: c.c.R(), verb: verb}
	return req
}

func (c *RESTClient) Get() *Request {
	return c.Verb("GET")
}

func (c *RESTClient) Post() *Request {
	return c.Verb("POST")
}

func (c *RESTClient) Put() *Request {
	return c.Verb("PUT")
}

func (c *RESTClient) Delete() *Request {
	return c.Verb("DELETE")
}

func (c *RESTClient) Head() *Request {
	return c.Verb("HEAD")
}

func (c *RESTClient) Cookies() []*http.Cookie {
	return c.c.Cookies
}
