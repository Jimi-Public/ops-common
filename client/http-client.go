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
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/Jimi-Public/ops-common/jwt"
	"github.com/Jimi-Public/ops-common/log"
)

var DefaultRequestBuilder = requestBuilder{}

type HttpClientBuilder interface {
	SetMethod(method string) HttpClientBuilder
	SetHeader(key, value string) HttpClientBuilder
	SetUrl(url string) HttpClientBuilder
	SetContext(ctx context.Context) HttpClientBuilder
	SetParams(params map[string]string) HttpClientBuilder
	SetBody(body string) HttpClientBuilder
	Build() (*http.Response, error)
}

type requestBuilder struct {
	headers map[string]string      // 请求头
	url     string                 // 请求路径
	method  string                 // 请求方法
	body    string                 // 请求体
	params  map[string]interface{} // get 传参
	ctx     context.Context        // 上下文
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

func (b *requestBuilder) SetParams(params map[string]string) HttpClientBuilder {
	b.params = make(map[string]interface{})
	for k, v := range params {
		b.params[k] = v
	}
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
	var err error
	c := &http.Client{
		Transport: http.DefaultTransport,
	}
	b.req, err = http.NewRequest(string(b.method), b.url, nil)
	b.req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
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

	switch strings.ToUpper(b.method) {
	case http.MethodGet:
		q := b.req.URL.Query()
		for k, v := range b.params {
			q.Add(k, v.(string))
		}
		b.req.URL.RawQuery = q.Encode()
	case http.MethodPost:
		b.req.Body = io.NopCloser(bytes.NewBuffer([]byte(b.body)))
	case http.MethodDelete:
		b.req.Body = io.NopCloser(bytes.NewBuffer([]byte(b.body)))
	default:
		return nil, errors.New("method not support")
	}
	return c.Do(b.req)
}
