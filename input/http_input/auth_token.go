package http_input

import (
	"errors"

	"github.com/zly-app/service/api"
)

// 验证token中间件
func AuthTokenMiddleware(token string) func(ctx *api.Context) error {
	return func(ctx *api.Context) error {
		if token == "" {
			return nil
		}
		if ctx.GetHeader("token") == token {
			return nil
		}
		return errors.New("Auth Token Failure")
	}
}
