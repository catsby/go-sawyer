package sawyer

import (
	"github.com/lostisland/go-sawyer/mediatype"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Request struct {
	Client    *http.Client
	ApiError  interface{}
	MediaType *mediatype.MediaType
	Query     url.Values
	*http.Request
}

func (c *Client) NewRequest(rawurl string, apierr interface{}) (*Request, error) {
	u, err := c.ResolveReferenceString(rawurl)
	if err != nil {
		return nil, err
	}

	httpreq, err := http.NewRequest(GetMethod, u, nil)
	for key, _ := range c.Header {
		httpreq.Header.Set(key, c.Header.Get(key))
	}
	return &Request{c.HttpClient, apierr, nil, httpreq.URL.Query(), httpreq}, err
}

func (r *Request) Do(method string, output interface{}) *Response {
	r.URL.RawQuery = r.Query.Encode()
	r.Method = method
	httpres, err := r.Client.Do(r.Request)
	if err != nil {
		return ResponseError(err)
	}

	mtype, err := mediaType(httpres)
	if err != nil {
		httpres.Body.Close()
		return ResponseError(err)
	}

	res := &Response{nil, mtype, UseApiError(httpres.StatusCode), false, httpres}
	if mtype != nil {
		res.decode(r.ApiError, output)
	}

	return res
}

func (r *Request) Head(output interface{}) *Response {
	return r.Do(HeadMethod, output)
}

func (r *Request) Get(output interface{}) *Response {
	return r.Do(GetMethod, output)
}

func (r *Request) Post(output interface{}) *Response {
	return r.Do(PostMethod, output)
}

func (r *Request) Put(output interface{}) *Response {
	return r.Do(PutMethod, output)
}

func (r *Request) Patch(output interface{}) *Response {
	return r.Do(PatchMethod, output)
}

func (r *Request) Delete(output interface{}) *Response {
	return r.Do(DeleteMethod, output)
}

func (r *Request) Options(output interface{}) *Response {
	return r.Do(OptionsMethod, output)
}

func (r *Request) SetBody(mtype *mediatype.MediaType, input interface{}) error {
	r.MediaType = mtype
	buf, err := mtype.Encode(input)
	if err != nil {
		return err
	}

	r.Header.Set(ctypeHeader, mtype.String())
	r.ContentLength = int64(buf.Len())
	r.Body = ioutil.NopCloser(buf)
	return nil
}

const (
	ctypeHeader   = "Content-Type"
	HeadMethod    = "HEAD"
	GetMethod     = "GET"
	PostMethod    = "POST"
	PutMethod     = "PUT"
	PatchMethod   = "PATCH"
	DeleteMethod  = "DELETE"
	OptionsMethod = "OPTIONS"
)
