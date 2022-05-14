package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type IHttpProxy interface {
	SetProxy(transport *http.Transport)
}

type HttpProxy struct {
	p  func(request *http.Request) (*url.URL, error)
	s5 ISocks5Proxy
}

func (h *HttpProxy) SetProxy(transport *http.Transport) {
	if h.s5 != nil {
		transport.DialContext = h.s5.DialContext
		return
	}

	transport.Proxy = h.p
}

/*创建一个http代理
  address 代理地址. 支持 http, https, socks5, socks5h. 示例: https://127.0.0.1:8080
  user 用户名
  passwd 密码
*/
func NewHttpProxy(address, user, passwd string) (IHttpProxy, error) {
	// 解析地址
	u, err := url.Parse(address)
	if err != nil {
		return nil, fmt.Errorf("address无法解析: %v", err)
	}

	scheme := strings.ToLower(u.Scheme)
	switch scheme {
	case "http", "https":
		if user != "" || passwd != "" {
			u.User = url.UserPassword(user, passwd)
		}
		p := func(request *http.Request) (*url.URL, error) {
			return u, nil
		}
		return &HttpProxy{p: p}, nil
	case "socks5", "socks5h":
		s5, err := NewSocks5Proxy(address, user, passwd)
		if err != nil {
			return nil, err
		}
		return &HttpProxy{s5: s5}, nil
	}
	return nil, fmt.Errorf("address的scheme不支持: %s")
}
