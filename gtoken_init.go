package token

import (
	"github.com/gogf/gf/frame/g"
)

var (
	LoginPath string
	LogoutPath string
	AuthPaths g.SliceStr
	AuthExcludePaths g.SliceStr
)

func New ()*GfToken{
	return &GfToken{
		ServerName:     g.Cfg().GetString("token.ServerName"),
		CacheMode:      g.Cfg().GetInt8("token.CacheMode"),
		CacheKey:       g.Cfg().GetString("token.CacheKey"),
		Timeout:        g.Cfg().GetInt("token.Timeout"),
		MaxRefresh:     g.Cfg().GetInt("token.MaxRefresh"),
		TokenDelimiter: g.Cfg().GetString("token.TokenDelimiter"),
		EncryptKey:     g.Cfg().GetBytes("token.EncryptKey"),
		AuthFailMsg:    g.Cfg().GetString("token.AuthFailMsg"),
		MultiLogin:     g.Config().GetBool("token.MultiLogin"),
		//GlobalMiddleware: g.Config().GetBool("token.GlobalMiddleware"),
		LoginPath: LoginPath,
		LogoutPath: LogoutPath,
		AuthPaths: AuthPaths,
		AuthExcludePaths: AuthExcludePaths,
	}
}
