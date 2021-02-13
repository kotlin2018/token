package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/kotlin2018/token"
)

var TestServerName  = "gtoken"

func main() {
	glog.Info("service start...")

	g.Cfg().SetPath("example/sample")
	s := g.Server(TestServerName)
	unAuth(s)
	auth(s)
	glog.Info("service finish.")
	s.Run()
}

// 不认证接口
func unAuth(s *ghttp.Server) {
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(token.CORS)
		// 调试路由
		group.ALL("/hello", func(r *ghttp.Request) {
			r.Response.WriteJson(token.Ok("hello"))
		})
	})
}

// 认证接口
func auth(s *ghttp.Server) {
	token.LoginPath = "/login"
	token.LogoutPath = "/user/logout"
	token.AuthExcludePaths = g.SliceStr{"/user/info", "/system/user/info"} // 不拦截路径 /user/info,/system/user/info,/system/user
	g2 := token.New()
	g2.LoginValidate = Login

	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(token.CORS)
		g2.Intercept(group)

		group.ALL("/system/user", func(r *ghttp.Request) {
			r.Response.WriteJson(token.Ok("system user"))
		})
		group.ALL("/user/data", func(r *ghttp.Request) {
			r.Response.WriteJson(g2.GetTokenData(r))
		})
		group.ALL("/user/info", func(r *ghttp.Request) {
			r.Response.WriteJson(token.Ok("user info"))
		})
		group.ALL("/system/user/info", func(r *ghttp.Request) {
			r.Response.WriteJson(token.Ok("system user info"))
		})
	})

	s.Group("/admin", func(group *ghttp.RouterGroup) {
		group.Middleware(token.CORS)
		g2.Intercept(group)

		group.ALL("/system/user", func(r *ghttp.Request) {
			r.Response.WriteJson(token.Ok("system user"))
		})
		group.ALL("/user/info", func(r *ghttp.Request) {
			r.Response.WriteJson(token.Ok("user info"))
		})
		group.ALL("/system/user/info", func(r *ghttp.Request) {
			r.Response.WriteJson(token.Ok("system user info"))
		})
	})
}

func Login(r *ghttp.Request) (string, interface{}) {
	username := r.GetString("username")
	passwd := r.GetString("passwd")

	if username == "" || passwd == "" {
		r.Response.WriteJson(token.Fail("账号或密码错误."))
		r.ExitAll()
	}

	return username, "1"
}

