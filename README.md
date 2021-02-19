# token
`token` 是基于`GoFrame`框架，对`gtoken`的简单修改和封装。

## 介绍
`gtoken` 是基于`GoFrame`框架的token插件，通过服务端验证方式实现token认证；已完全可以支撑线上token认证，通过Redis支持集群模式；使用简单，大家可以放心使用；

`token源码地址`
* github地址：https://github.com/kotlin2018/token
* gitee地址：https://gitee.com/KotlinToGo/token

`gtoken源码地址`
* github地址：https://github.com/goflyfox/gtoken
* gitee地址：https://gitee.com/goflyfox/gtoken

## gtoken优势
1. gtoken支撑单点应用使用内存存储，也支持集群使用redis存储；完全适用于企业生产级使用；
2. 有效的避免了jwt服务端无法退出问题；
3. 解决jwt无法作废已颁布的令牌，只能等到令牌过期问题；
4. 通过用户扩展信息存储在服务端，有效规避了jwt携带大量用户扩展信息导致降低传输效率问题；
5. 有效避免jwt需要客户端实现续签功能，增加客户端复杂度；支持服务端自动续期，客户端不需要关心续签逻辑；

## 特性说明

1. 支持token认证，不强依赖于session和cookie，适用jwt和session认证所有场景；
2. 支持单机gcache和集群gredis模式；
```
# 缓存模式 1 gcache 2 gredis
CacheMode = 2
```

2. 支持服务端缓存自动续期功能
```
// 注：通过MaxRefresh，默认当用户第五天访问时，自动续期
// 超时时间 默认10天
Timeout int
// 缓存刷新时间 默认为超时时间的一半
MaxRefresh int
```
4. 支持分组拦截、全局拦截、深度路径拦截，便于根据个人需求定制拦截器；**建议使用分组拦截方式；**
5. 框架使用简单，只需要设置登录验证方法以及登录、登出路径即可；
6. 在`gtoken v1.4.0`版本开始支持分组中间件方式实现，但依然兼容全局和深度中间件实现方式；

## 安装教程

* gopath模式: `go get github.com/kotlin2018/token`
* 或者 使用go.mod添加 :`require github.com/kotlin2018/token latest`

## 分组中间件实现

GoFrame官方推荐使用Group方式实现路由和中间件；

### 使用说明

推荐使用分组方式实现

```go
package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/kotlin2018/token"
)

func main() {
	glog.Info("service start...")
	g.Cfg().SetPath("example/sample1")
	s := g.Server(TestServerName)
	group(s)
	glog.Info("service finish...")
	s.Run()
}

// 分组拦截
func group(s *ghttp.Server) {
	g2 := token.GfToken{}
	g2.LoginPath 		= "/login"
	g2.LogoutPath 		= "/user/logout"
	g2.AuthExcludePaths 	= g.SliceStr{"/user/info", "/system/user/*"} // 不拦截路径  /user/info,/system/user/info,/system/user
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
```

登录方法实现，通过username返回空或者r.ExitAll()\r.Exit()处理认证失败；

````go
func Login(r *ghttp.Request) (string, interface{}) {
    username := r.GetString("username")
    passwd := r.GetString("passwd")
    
    if username == "" || passwd == "" {
        r.Response.WriteJson(token.Fail("账号或密码错误."))
        r.ExitAll()
    }
    return username, "1"
}
````

通过`gtoken.GetTokenData(r)`获取登录信息，例如: 获取用户名...

### 路径拦截规则

```go
    AuthExcludePaths: g.SliceStr{"/user/info", "/system/user/*"}, // 不拦截路径  /user/info,/system/user/info,/system/user,
```

1. 分组中间件实现，不需要设置`AuthPaths`认证路径，设置也没有作用，**需要认证路径为该分组下所有路由**；
2. 使用分组拦截的是通过GoFrame的`group.Middleware(authMiddleware)`方法，对该分组下的所有路由进行拦截；
3. 对登录接口路径`loginPath`和登出接口路径`logoutPath`做拦截认证放行，登出放行是为了避免认证过期无法登出情况；
4. 严格按照GoFrame分组中间件拦截优先级；如果使用跨域中间件，建议放在跨域中间件之后；
5. 如果配置`AuthExcludePaths`路径，会将配置的不拦截路径排除；

### 逻辑测试

参考sample项目，先运行main.go，然后可运行api_test.go进行测试并查看结果；验证逻辑说明：

1. 访问用户信息，提示未携带token
2. 调用登录后，携带token访问正常
3. 调用登出提示成功
4. 携带之前token访问，提示未登录

```json
=== RUN   TestAdminSystemUser
    api_admin_test.go:22: 1. not login and visit user
    api_admin_test.go:29: {"code":-401,"msg":"请求错误或登录超时","data":""}
    api_admin_test.go:42: 2. execute login and visit user
    api_admin_test.go:45: {"code":0,"msg":"success","data":"system user"}
    api_admin_test.go:51: 3. execute logout
    api_admin_test.go:54: {"code":0,"msg":"success","data":"Logout success"}
    api_admin_test.go:60: 4. visit user
    api_admin_test.go:65: {"code":-401,"msg":"请求错误或登录超时","data":""}
```

## 全局中间件实现

### 使用说明

只需要配置登录路径、登出路径、拦截路径以及登录校验实现即可

```go
package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/kotlin2018/token"
)

func main() {
	glog.Info("service start...")
	g.Cfg().SetPath("example/sample1")
	s := g.Server(TestServerName)
	global(s)
	glog.Info("service finish...")
	s.Run()
}

// 全局拦截
func global(s *ghttp.Server){
	// 默认从配置文件中读取配置项
	g2 := token.GfToken{}
	g2.LoginPath 		= "/login"
	g2.LogoutPath 		= "/user/logout"
	g2.AuthPaths  		= g.SliceStr{"/user", "/system"}                // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
	g2.AuthExcludePaths     = g.SliceStr{"/user/info", "/system/user/info"} // 不拦截路径  /user/info,/system/user/info,/system/user,
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
````

`登录方法实现，通过username返回空或者r.ExitAll()\r.Exit()处理认证失败；`
```go
func Login(r *ghttp.Request) (string, interface{}) {
	username := r.GetString("username")
	passwd := r.GetString("passwd")
	
	// TODO 进行登录校验
	if username == "" || passwd == "" {
		r.Response.WriteJson(token.Fail("账号或密码错误."))
		r.ExitAll()
	}
	return username, "1"
}
```

通过`gtoken.GetTokenData(r)`获取登录信息

### 路径拦截规则
```go
    AuthPaths:        g.SliceStr{"/user", "/system"},             // 这里是按照前缀拦截，拦截/user /user/list /user/add ...
    AuthExcludePaths: g.SliceStr{"/user/info", "/system/user/*"}, // 不拦截路径  /user/info,/system/user/info,/system/user,
```

1. `GlobalMiddleware:true`全局拦截的是通过GF的`BindMiddleware`方法创建拦截`/*`
2. `GlobalMiddleware:false`路径拦截的是通过GF的`BindMiddleware`方法创建拦截`/user*和/system/*`
3. 按照中间件优先级路径拦截优先级很高；如果先实现部分中间件在认证前处理需要切换成全局拦截器，严格按照注册顺序即可；
4. 程序先处理认证路径，如果满足；再排除不拦截路径；
5. 如果只想用排除路径功能，将拦截路径设置为`/*`即可；

### 逻辑测试

参考sample1项目，先运行main.go，然后可运行api_test.go进行测试并查看结果；验证逻辑说明：

1. 访问用户信息，提示未携带token
2. 调用登录后，携带token访问正常
3. 调用登出提示成功
4. 携带之前token访问，提示未登录

```json
=== RUN   TestSystemUser
    api_test.go:43: 1. not login and visit user
    api_test.go:50: {"code":-401,"msg":"请求错误或登录超时","data":""}
    api_test.go:63: 2. execute login and visit user
    api_test.go:66: {"code":0,"msg":"success","data":"system user"}
    api_test.go:72: 3. execute logout
    api_test.go:75: {"code":0,"msg":"success","data":"Logout success"}
    api_test.go:81: 4. visit user
    api_test.go:86: {"code":-401,"msg":"请求错误或登录超时","data":""}
```
## 返回码及配置项

1. 正常操作成功返回0
2. 未登录访问需要登录资源返回401
3. 程序异常返回-99，如编解码错误等

```go
SUCCESS      = 0  // 正常
FAIL         = -1  // 失败
ERROR        = -99  // 异常
UNAUTHORIZED = -401  // 未认证
```
### 配置项说明

具体可参考`GfToken`结构体，字段解释如下：

| 名称             | 配置字段         | 说明                                                    | 分组中间件 | 全局中间件 |
| ---------------- | ---------------- | ------------------------------------------------------- | ---------- | ---------- |
| 服务名           | ServerName       | 默认空即可                                              | 支持       | 支持       |
| 缓存模式         | CacheMode        | 1 gcache 2 gredis 默认1                                 | 支持       | 支持       |
| 缓存key          | CacheKey         | 默认缓存前缀`GToken:`                                   | 支持       | 支持       |
| 超时时间         | Timeout          | 默认10天（毫秒）                                        | 支持       | 支持       |
| 缓存刷新时间     | MaxRefresh       | 默认为超时时间的一半（毫秒）                            | 支持       | 支持       |
| Token分隔符      | TokenDelimiter   | 默认`_`                                                 | 支持       | 支持       |
| Token加密key     | EncryptKey       | 默认`12345678912345678912345678912345`                  | 支持       | 支持       |
| 认证失败提示     | AuthFailMsg      | 默认`请求错误或登录超时`                                | 支持       | 支持       |
| 是否支持多端登录 | MultiLogin       | 默认false                                               | 支持       | 支持       |
| 中间件类型       | MiddlewareType   | 1、Group 2、Bind 3 、Global；<br>使用分组模式不需要设置 | 支持       | 支持       |
| 登录路径         | LoginPath        | 登录接口路径                                            | 支持       | 支持       |
| 登录验证方法     | LoginValidate  | 登录验证需要用户实现方法                                | 支持       | 支持       |
| 登录返回方法     | LoginResponse   | 登录完成后调用                                          | 支持       | 支持       |
| 登出地址         | LogoutPath       | 登出接口路径                                            | 支持       | 支持       |
| 登出验证方法     | LogoutValidate | 登出接口前调用                                          | 支持       | 支持       |
| 登出返回方法     | LogoutResponse  | 登出接口完成后调用                                      | 支持       | 支持       |
| 拦截地址         | AuthPaths        | 此路径列表进行认证                                      | 不需要     | 支持       |
| 拦截排除地址     | AuthExcludePaths | 此路径列表不进行认证                                    | 支持       | 支持       |
| 认证验证方法     | Authenticator   | 拦截认证前调用                                        | 支持       | 支持       |
| 认证返回方法     | AuthResponse    | 拦截认证完成后调用                                      | 支持       | 支持       |

## 感谢

1. gf框架 [https://github.com/gogf/gf](https://github.com/gogf/gf) 

