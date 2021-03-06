package http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type ResponseBuilderValues struct {
	StatusCode int
	Status     string
	Request    *http.Request
	Header     http.Header
	Body       io.ReadCloser
	Version    string
}

type ResponseBuilder struct {
	Values ResponseBuilderValues
}

func NewResponseBuilder(req *http.Request, version string) *ResponseBuilder {
	return &ResponseBuilder{
		Values: ResponseBuilderValues{
			StatusCode: 200,
			Status:     "200 OK",
			Request:    req,
			Header:     CreateHeader(version),
			Body:       ioutil.NopCloser(&bytes.Buffer{}),
			Version:    "",
		},
	}
}

func (b *ResponseBuilder) StatusCode(statusCode int) *ResponseBuilder {
	b.Values.StatusCode = statusCode
	return b
}

func (b *ResponseBuilder) Status(status string) *ResponseBuilder {
	b.Values.Status = status
	return b
}

func (b *ResponseBuilder) AddHeader(key string, value string) *ResponseBuilder {
	b.Values.Header.Add(key, value)
	return b
}

func (b *ResponseBuilder) BodyString(body string) *ResponseBuilder {
	return b.BodyBytes([]byte(body))
}

func (b *ResponseBuilder) BodyBytes(body []byte) *ResponseBuilder {
	return b.Body(ioutil.NopCloser(bytes.NewReader(body)))
}

func (b *ResponseBuilder) Body(body io.ReadCloser) *ResponseBuilder {
	b.Values.Body = body
	return b
}

func (b *ResponseBuilder) Version(version string) *ResponseBuilder {
	b.Values.Version = version
	return b
}

func (b *ResponseBuilder) Build() *http.Response {
	req := b.Values.Request
	return &http.Response{
		Proto:      req.Proto,
		ProtoMajor: req.ProtoMajor,
		ProtoMinor: req.ProtoMinor,
		StatusCode: b.Values.StatusCode,
		Status:     b.Values.Status,
		Header:     b.Values.Header,
		Body:       b.Values.Body,
		Request:    req,
	}
}

func CreateHeader(version string) http.Header {
	header := make(http.Header)
	proxyName := fmt.Sprintf("Granblue Proxy %s", version)
	header.Add("X-Proxy-Server", proxyName)
	return header
}
