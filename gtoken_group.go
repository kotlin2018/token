package token

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

// Intercept 一个 GfToken 实例拦截一个 *ghttp.RouterGroup (路由组)下面所有子路由
//
// 这个被拦截的路由组能够使用GfToken一些功能，例如: 路径拦截，登陆验证...
func (m *GfToken) Intercept(group *ghttp.RouterGroup) bool {
	if !m.InitConfig() {
		return false
	}
	// 设置为Group模式
	m.MiddlewareType = MiddlewareTypeGroup
	glog.Info("[GToken][params:" + m.String() + "]start... ")

	// 缓存模式
	if m.CacheMode > CacheModeRedis {
		glog.Error("[GToken]CacheMode set error")
		return false
	}
	// 登录
	if m.LoginPath == "" || m.LoginValidate == nil {
		glog.Error("[GToken]LoginPath or LoginBeforeFunc not set")
		return false
	}
	// 登出
	if m.LogoutPath == "" {
		glog.Error("[GToken]LogoutPath or LogoutFunc not set")
		return false
	}

	group.Middleware(m.authMiddleware)
	group.ALL(m.LoginPath, m.Login)
	group.ALL(m.LogoutPath, m.Logout)

	return true
}
