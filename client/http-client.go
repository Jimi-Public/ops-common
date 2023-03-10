/*
@Project: ops-common
@Author:  WangChaoQun
@Date:    2023/2/13
@IDE:     GoLand
@File:    http-client.go
统一构造 http request 请求
*/

package client

import (
	"bytes"
	"context"
	"net/http"

	"github.com/Jimi-Public/ops-common/jwt"
	"github.com/Jimi-Public/ops-common/log"
)

var DefaultRequestBuilder = &requestBuilder{}

type HttpClientBuilder interface {
	SetMethod(method string) HttpClientBuilder
	SetHeader(key, value string) HttpClientBuilder
	SetUrl(url string) HttpClientBuilder
	SetContext(ctx context.Context) HttpClientBuilder
	SetRequest(req *http.Request) HttpClientBuilder
	Build() (*http.Response, error)
}

type requestBuilder struct {
	headers map[string]string // 请求头
	url     string            // 请求路径
	method  string            // 请求方法
	body    string            // 请求体
	ctx     context.Context   // 上下文
	req     *http.Request
}

// SetMethod 设置请求方法
func (b *requestBuilder) SetMethod(method string) HttpClientBuilder {
	b.method = method
	return b
}

// SetRequest 自定义Request
func (b *requestBuilder) SetRequest(req *http.Request) HttpClientBuilder {
	b.req = req
	return b
}

// SetHeader 设置请求头
func (b *requestBuilder) SetHeader(key, value string) HttpClientBuilder {
	if b.headers == nil {
		b.headers = make(map[string]string)
	}
	b.headers[key] = value
	return b
}

// SetUrl 设置Url
func (b *requestBuilder) SetUrl(url string) HttpClientBuilder {
	b.url = url
	return b
}

// SetBody 设置Body
func (b *requestBuilder) SetBody(body string) HttpClientBuilder {
	b.body = body
	return b
}

// SetContext context
func (b *requestBuilder) SetContext(ctx context.Context) HttpClientBuilder {
	b.ctx = ctx
	return b
}

func (b *requestBuilder) Build() (*http.Response, error) {
	c := &http.Client{
		Transport: http.DefaultTransport,
	}
	if b.req == nil {
		var err error
		b.req, err = http.NewRequest(b.method, b.url, bytes.NewBuffer([]byte(b.body)))
		if err != nil {
			return nil, err
		}
	}

	if b.ctx != nil {
		// Context 中获取Token 透传
		if token := b.ctx.Value(jwt.AuthHeader); token != nil {
			b.SetHeader(jwt.AuthHeader, token.(string))
		}
		//  Context 中获取Trace-id 透传
		if id := b.ctx.Value(log.TraceName); id != nil {
			b.SetHeader(log.TraceName, id.(string))
		}
	}
	for k, v := range b.headers {
		b.req.Header.Add(k, v)
	}
	return c.Do(b.req)
}
