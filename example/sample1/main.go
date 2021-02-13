package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/kotlin2018/token"
)

var TestServerName = "gToken"

func main() {
	glog.Info("service start...")

	g.Cfg().SetPath("example/sample1")
	s := g.Server(TestServerName)

	//group(s)
	global(s)

	glog.Info("service finish.")
	s.Run()
}

// 分组拦截
func group(s *ghttp.Server) {
	token.LoginPath 		= "/login"
	token.LogoutPath 		= "/user/logout"
	token.AuthExcludePaths 	= g.SliceStr{"/user/info", "/system/user/*"} // 不拦截路径  /user/info,/system/user/info,/system/user
	g2 := token.New()
	g2.LoginValidate = Login

	s.Group("/base", func(group *ghttp.RouterGroup) {
		group.Middleware(token.CORS)
		g2.Intercept(group)

		group.ALL("/system/user", func(r *ghttp.Request) {
			r.Response.WriteJson(token.Ok("system user"))
		})
	})

	s.Group("/", func(group *ghttp.RouterGroup) { // 这个 "/system/user/list" 在 token.AuthExcludePaths里，所以不拦截。
		group.Middleware(token.CORS)
		g2.Intercept(group)
		group.ALL("/system/user/list", func(r *ghttp.Request) {
			r.Response.WriteJson(token.Ok("system user list"))
		})
	})
}

// 全局拦截
func global(s *ghttp.Server){
	token.LoginPath 		= "/login"
	token.LogoutPath 		= "/user/logout"
	token.AuthPaths  		= g.SliceStr{"/user", "/system"}                // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
	token.AuthExcludePaths 	= g.SliceStr{"/user/info", "/system/user/info"} // 不拦截路径  /user/info,/system/user/info,/system/user,
	g2 := token.New()
	g2.LoginValidate = Login

	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(token.CORS)
		// 调试路由
		group.ALL("/hello", func(r *ghttp.Request) {  // 这个不拦截，拦截规则没有约束它 (因为它既不在: 拦截路径列表里,也不在: 不拦截路径列表里)。
			r.Response.WriteJson(token.Ok("hello"))
		})
		group.ALL("/system/user", func(r *ghttp.Request) { // 这个拦截，因为它在 :拦截路径列表里
			r.Response.WriteJson(token.Ok("system user"))
		})
		group.ALL("/user/info", func(r *ghttp.Request) { // 这个不拦截，因为它在 :不拦截路径列表里
			r.Response.WriteJson(token.Ok("user info"))
		})
		group.ALL("/system/user/info", func(r *ghttp.Request) { // 这个不拦截，因为它在 :不拦截路径列表里
			r.Response.WriteJson(token.Ok("system user info"))
		})
	})

	s.Group("/base", func(group *ghttp.RouterGroup) { // 这个不拦截，因为拦截规则没有约束它。
		group.Middleware(token.CORS)
		group.ALL("/system/user", func(r *ghttp.Request) {
			r.Response.WriteJson(token.Ok("base system user"))
		})
	})
	g2.Start()
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

