# 开发 web 服务程序
## 概述
开发简单 web 服务程序 cloudgo，了解 web 服务器工作原理。

## 任务目标
- 熟悉 go 服务器工作原理
- 基于现有 web 库，编写一个简单 web 应用类似 cloudgo。
- 使用 curl 工具访问 web 程序
- 对 web 执行压力测试
## 任务要求
- 基本要求
	* 编程 web 服务程序 类似 cloudgo 应用。
		+ 支持静态文件服务
		+ 支持简单 js 访问
		+ 提交表单，并输出一个表格（必须使用模板）
	* 使用 curl 测试，将测试结果写入readme
	* 使用 ab 测试，将测试结果写入 README。并解释重要参数。
- 扩展要求
选择以下一个或多个任务，以博客的形式提交。
	* 通过源码分析、解释一些关键功能实现
	* 选择简单的库，如 mux 等，通过源码分析、解释它是如何实现扩展的原理，包括一些 golang 程序设计技巧。
## 注意事项
本博客由于是在csdn上进行撰写，所以含有图片水印，csdn账号名：qq_43283265
## 实验过程
### 应用功能部分
#### 服务器开启
利用`server := service.NewServer()`创建新的连接，这个函数将在下面的具体实现部分给出

```go
package main

import (
	"os"

	service "github.com/yilin0041/service-computing/cloudgo"

	flag "github.com/spf13/pflag"
)

const (
	PORT string = "8080"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = PORT
	}
	pPort := flag.StringP("port", "p", PORT, "PORT for httpd listening")
	flag.Parse()
	if len(*pPort) != 0 {
		port = *pPort
	}
	server := service.NewServer()
	server.Run(":" + port)
}

```

#### 静态文件服务功能
本部分参考老师给出的代码，我们可以知道，在开发、测试阶段部署 Nginx 服务器太麻烦。net/http 库以提供了现成的支持，仅使用一个函数 `func FileServer(root FileSystem) Handler` 就可以使用上文件服务。具体代码如下：

```go
package service

import (
	"fmt"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {

	// formatter := render.New(render.Options{
	// 	IndentJSON: true,
	// })
	formatter := render.New(render.Options{
		Directory:  "templates",
		Extensions: []string{".html"},
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()

	initRoutes(mx, formatter)

	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	webRoot := os.Getenv("WEBROOT")
	if len(webRoot) == 0 {
		if root, err := os.Getwd(); err != nil {
			panic("Could not retrive working directory")
		} else {
			webRoot = root
			fmt.Println(root)
		}
	}
	mx.PathPrefix("/").Handler(http.FileServer(http.Dir(webRoot + "/assets/")))

}
```
首先我们需要在服务器上创建目录，以存放静态内容。例如：
>assets(静态文件虚拟根目录)
  |-- js
  |-- images
  +-- css

然后，我们运行`go run main.go`，结果如下所示
![file](https://img-blog.csdnimg.cn/20201123183225358.png)

然后我们任意打开几个文件进心检查：

- 打开1.txt

![test1](https://img-blog.csdnimg.cn/2020112318331612.png)

- 打开一张图片

![test2](https://img-blog.csdnimg.cn/20201123183357735.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)

经测试，静态文件功能可以访问

#### 简单js访问功能
随着 web 页面技术的进步，页面中大量使用 javascript。 所有我们添加一个服务：apitest.go

```go
package service

import (
	"net/http"

	"github.com/unrolled/render"
)

func apiTestHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct {
			ID      string `json:"id"`
			Content string `json:"content"`
		}{ID: "8675309", Content: "Hello from Go!"})
	}
}
```
用 curl 测试输出
![test](https://img-blog.csdnimg.cn/20201123183752235.png)
在`assets`中创建一个文件`jsonindex.html`，使用简单的代码进行测试，代码如下：

```html
<html>
<head>
  <link rel="stylesheet" href="css/main.css"/>
  <script src="http://code.jquery.com/jquery-latest.js"></script>
  <script src="js/hello.js"></script>
</head>
<body>
  <img src="images/cng.png" height="48" width="48"/>
  Sample Go Web Application!!
      <div>
          <p class="greeting-id">The ID is </p>
          <p class="greeting-content">The content is </p>
      </div>
</body>
</html>
```
这里使用的`hello.js`为

```js
$(document).ready(function() {
    $.ajax({
        url: "/api/test"
    }).then(function(data) {
       $('.greeting-id').append(data.id);
       $('.greeting-content').append(data.content);
    });
});
```
最后在server中添加路由

```go
mx.HandleFunc("/api/test", apiTestHandler(formatter)).Methods("GET")
```
然后运行打开`jsonindex.html`，结果如下：
![jsontest](https://img-blog.csdnimg.cn/2020112318442541.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)

可以看到后台的访问有js的接口，证明数据是通过js传入的

![jstest](https://img-blog.csdnimg.cn/20201123184623634.png)

#### 使用模板完成提交表单输出表格功能

这里需要使用模板完成，所以我们要新建一个`templates`文件夹，然后在服务器设置的时候，我们已经写明了在这个文件夹下读取html模板，所以只需要将模板放入这里即可。
先按照参考文档中的例子进行一个简单的index页面制作，首先编辑`index.html`

```html
<html>
<head>
  <link rel="stylesheet" href="css/main.css"/>
</head>
<body>
  <img src="images/cng.png" height="48" width="48"/>
  Sample Go Web Application!!
      <div>
          <p class="greeting-id">The ID is {{.ID}}</p>
          <p class="greeting-content">The content is {{.Content}}</p>
      </div>
</body>
</html>
```
新建一个`homeHandler.go`来进行返回

```go
package service

import (
	"net/http"

	"github.com/unrolled/render"
)

func homeHandler(formatter *render.Render) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		formatter.HTML(w, http.StatusOK, "index", struct {
			ID      string `json:"id"`
			Content string `json:"content"`
		}{ID: "8675309", Content: "Hello from Go!"})
	}
}

```
我们使用 formatter 的 HTML 直接将数据注入模板，并输出到浏览器。这里创建一个路由

```go
mx.HandleFunc("/index", homeHandler(formatter))
```
下面进行一下测试，发现可以正常运行：
![test](https://img-blog.csdnimg.cn/20201123185716851.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)

发现功能正常

下面来实现表单的提交

首先来写一个提交用的html模板`login.html`

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Login</title>
    <link rel="stylesheet" type="text/css" href="css/login.css">
</head>
<body>

<div id="content">
    <div class="login-header">
        <img src="images/cng.png">
    </div>
    <form action="/user" method="post">
        <div class="login-input-box">
            <span class="icon icon-user"></span>
            <input type="text" name="username" placeholder="Please enter your username">
        </div>
        <div class="login-input-box">
            <span class="icon icon-password"></span>
            <input type="password" name="password" placeholder="Please enter your password">
		</div>
		<div class="remember-box">
			<label>
				<input type="checkbox"> Remember Me
			</label>
		</div>
		<div class="login-button-box">
			<button>Login</button>
		</div>
    </form>
</div>
</body>
</html>
```
这里css不是重点，故不再赘述。

新建一个`login.go`来进行返回，这里设定账户为`sunylin`，密码为`123`，如果正确则进入`user.html`模板，如果错误则进入`error.html`模板

```go
func login(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			formatter.HTML(w, http.StatusOK, "login", struct{}{})
		} else {
			req.ParseForm()
			if req.Form["username"][0] == "sunylin" && req.Form["password"][0] == "123" {
				formatter.HTML(w, http.StatusOK, "user", struct {
					Username string
					Password string
				}{
					Username: req.Form["username"][0],
					Password: req.Form["password"][0],
				})
			} else {
				formatter.HTML(w, http.StatusOK, "error", struct{}{})
			}
		}
	}
}
```
user模板显示正确的账号和密码，代码如下

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Login</title>
    <link rel="stylesheet" type="text/css" href="css/user.css">
</head>
<body>

<div id="content">
    <div class="login-header">
        <img src="images/cng.png">
    </div>
        <div class="login-input-box">
            <label>Username:{{.Username}}</label>
        </div>
        <div class="login-input-box">
            <label>Password:{{.Password}}</label>
		</div>
</div>
</body>
</html>
```
error模板则显示报错，代码如下

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Login</title>
    <link rel="stylesheet" type="text/css" href="css/error.css">
</head>
<body>

<div id="content">
    <div class="login-header">
        <img src="images/cng.png">
    </div>
        <div class="login-input-box">
            <label>Username or Password Error</label>
        </div>
</div>
</body>
</html>
```
最后我们要添加路由设置

```go
mx.HandleFunc("/login", login(formatter))
mx.HandleFunc("/user", login(formatter))
mx.HandleFunc("/error", login(formatter))
```
测试结果如下：

![test](https://img-blog.csdnimg.cn/2020112319142680.gif)

后台截图如下：

![test](https://img-blog.csdnimg.cn/20201123191607995.png)

### 测试部分
#### curl测试

##### 文件系统

![test](https://img-blog.csdnimg.cn/2020112319215513.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)

##### 访问文件系统中的`jsonindex.html`
![test](https://img-blog.csdnimg.cn/20201123193258435.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)

##### index页面
![test](https://img-blog.csdnimg.cn/20201123192950507.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)

##### 登录界面

![test](https://img-blog.csdnimg.cn/20201123192631273.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)

##### 用户界面

![test](https://img-blog.csdnimg.cn/20201123192738579.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)

##### 错误页

![test](https://img-blog.csdnimg.cn/20201123192822841.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)
#### ab测试
在文件系统的测试中有结果的详细解释
##### 文件系统

首先对文件系统做ab测试

```bash
ab -n 100000 -c 100 http://localhost:8080/
```
这里发送10万请求，100并发
测试结果如下：
![test](https://img-blog.csdnimg.cn/20201123194538137.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)
结果解释：
- Concurrency Level:并发量
- Time taken for tests: 测试用时
- Complete requests:完成请求次数
- Failed requests:失败请求次数
- Total transferred:总传输字节数
- HTML transferred:html传输字节数
- Requests per second: 每秒平均请求次数
- Time per request:每个请求平均用时（并发）
- Time per request: 并发中每个请求平均用时
- Transfer rate:带宽速度
- connection time表格，用于解释最小、最大、中等、平均值
- 最后是一个响应时间分布表


##### 访问文件系统中的`jsonindex.html`
```bash
ab -n 10000 -c 100 http://localhost:8080/jsonindex.html
```
![test](https://img-blog.csdnimg.cn/20201123195934402.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)
##### index页面
```bash
ab -n 10000 -c 100 http://localhost:8080/index
```
![test](https://img-blog.csdnimg.cn/20201123200057786.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)

##### 登录界面

```bash
ab -n 10000 -c 100 http://localhost:8080/login
```
![test](https://img-blog.csdnimg.cn/2020112319524220.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)
##### 用户页面
```bash
ab -n 10000 -c 100 http://localhost:8080/user
```
![test](https://img-blog.csdnimg.cn/20201123200138513.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)

##### 错误页
```bash
ab -n 10000 -c 100 http://localhost:8080/error
```
![test](https://img-blog.csdnimg.cn/20201123200217855.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjgzMjY1,size_16,color_FFFFFF,t_70)
### 源代码阅读
为了更好的阅读，参考了老师给出的参考书：[Go如何使得Web工作](https://gitee.com/astaxie/build-web-application-with-golang/blob/master/zh/03.3.md)
首先，我们一定要先了解是Go实现Web服务的工作模式的流程图，如下所示：
![web](https://img-blog.csdnimg.cn/img_convert/cebdee002ad171a90ac9e6287f0b774c.png)
书中告诉我们，http包执行流程大致为
1. 创建Listen Socket, 监听指定的端口, 等待客户端请求到来。

2. Listen Socket接受客户端的请求, 得到Client Socket, 接下来通过Client Socket与客户端通信。

3. 处理客户端的请求, 首先从Client Socket读取HTTP请求的协议头, 如果是POST方法, 还可能要读取客户端提交的数据, 然后交给相应的handler处理请求, handler处理完毕准备好客户端需要的数据, 通过Client Socket写给客户端。 

如何实现的呢？？下面开始阅读源码来寻找答案
首先找到入口函数`ListenAndServe`，初始化一个`server`对象，然后调用了`net.Listen("tcp", addr)`，也就是底层用TCP协议搭建了一个服务，然后监控我们设置的端口。

```go
func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}
```

```go
func (srv *Server) ListenAndServe() error {
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}
```
下面是`Serve`函数，这个函数就是处理接收客户端的请求信息。`testHookServerServe`执行net/http库默认的测试函数。`tempDelay` 设置一个最长相应时间，然后调用`setupHTTP2_Serve()`
`setupHTTP2_Serve`设置http2，如下
```go
func (srv *Server) setupHTTP2_Serve() error {
	srv.nextProtoOnce.Do(srv.onceSetNextProtoDefaults_Serve)
	return srv.nextProtoErr
}
```
然后`trackListener`设置track日志，`baseCtx`是Server一个监听的根Context
下面进入for循环
`Accept`获得了一个net.Conn连接对象，使用`srv.newConn(rw)`方法创建一个`http.conn`连接。

> http.conn连接就是http连接

设置连接状态用于连接复用，然后c.serve处理这个http连接。


最后附上源码：
```go
func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()
	if fn := testHookServerServe; fn != nil {
		fn(srv, l)
	}
	var tempDelay time.Duration // how long to sleep on accept failure

	if err := srv.setupHTTP2_Serve(); err != nil {
		return err
	}

	srv.trackListener(l, true)
	defer srv.trackListener(l, false)

	baseCtx := context.Background() // base is always background, per Issue 16220
	ctx := context.WithValue(baseCtx, ServerContextKey, srv)
	for {
		rw, e := l.Accept()
		if e != nil {
			select {
			case <-srv.getDoneChan():
				return ErrServerClosed
			default:
			}
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				srv.logf("http: Accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return e
		}
		tempDelay = 0
		c := srv.newConn(rw)
		c.setState(c.rwc, StateNew) // before Serve can return
		go c.serve(ctx)
	}
}

```


net.conn.serve处理一个http连接,由于代码很长，这里只看重点部分
使用Hijack方法获取了tcp连接，然后自定义处理，c.rwc就是net.Conn的连接。

```go
		if !c.hijacked() {
			c.close()
			c.setState(c.rwc, StateClosed)
		}
```
for循环调用`readRequest`读取一个请求，并创建`ResponseWriter`对象。

```go
	for {
		w, err := c.readRequest(ctx)
		if c.r.remain != c.server.initialReadLimitSize() {
			// If we read any bytes off the wire, we're active.
			c.setState(c.rwc, StateActive)
		}
```

创建一个serverHandler处理当前的请求rw，serverHandler就检测Server是否设置了默认处理者，和响应Option方法。
```go
		serverHandler{c.server}.ServeHTTP(w, w.req)
		w.cancelCtx()
		if c.hijacked() {
			return
		}
```
`readRequest`方法就是根据连接创建`http.Request`和`http.ResponseWriter`两个对象供`http.Handler`接口使用，处理一个请求,创建过程如下：
```go
	w = &response{
		conn:          c,
		cancelCtx:     cancelCtx,
		req:           req,
		reqBody:       req.Body,
		handlerHeader: make(Header),
		contentLength: -1,
		closeNotifyCh: make(chan bool, 1),

		// We populate these ahead of time so we're not
		// reading from req.Header after their Handler starts
		// and maybe mutates it (Issue 14940)
		wants10KeepAlive: req.wantsHttp10KeepAlive(),
		wantsClose:       req.wantsClose(),
	}
```
最后来看一下`Handler`函数，`Handler`返回用于给定请求的处理程序，查询`Method`、`Host`和`URL`.它返回一个非nil处理程序。如果路径不是其规范形式，则处理程序将是内部生成的重定向到规范路径的处理程序。如果主机包含端口，则在匹配处理程序时将忽略该端口。如果没有应用于请求的注册处理程序，处理程序返回“找不到页面”处理程序和空模式。
```go
func (mux *ServeMux) Handler(r *Request) (h Handler, pattern string) {

	// CONNECT requests are not canonicalized.
	if r.Method == "CONNECT" {
		// If r.URL.Path is /tree and its handler is not registered,
		// the /tree -> /tree/ redirect applies to CONNECT requests
		// but the path canonicalization does not.
		if u, ok := mux.redirectToPathSlash(r.URL.Host, r.URL.Path, r.URL); ok {
			return RedirectHandler(u.String(), StatusMovedPermanently), u.Path
		}

		return mux.handler(r.Host, r.URL.Path)
	}

	// All other requests have any port stripped and path cleaned
	// before passing to mux.handler.
	host := stripHostPort(r.Host)
	path := cleanPath(r.URL.Path)

	// If the given path is /tree and its handler is not registered,
	// redirect for /tree/.
	if u, ok := mux.redirectToPathSlash(host, path, r.URL); ok {
		return RedirectHandler(u.String(), StatusMovedPermanently), u.Path
	}

	if path != r.URL.Path {
		_, pattern = mux.handler(host, path)
		url := *r.URL
		url.Path = path
		return RedirectHandler(url.String(), StatusMovedPermanently), pattern
	}

	return mux.handler(host, r.URL.Path)
}
```
最后，利用参考资料中的总结再来梳理一下：
调用`http.ListenAndServe`按顺序做了几件事情：
- 实例化Server
- 调用Server的ListenAndServe()
- 调用net.Listen("tcp", addr)监听端口
- 启动一个for循环，在循环体中Accept请求
- 对每个请求实例化一个Conn，并且开启一个goroutine为这个请求进行服务go c.serve()
- 读取每个请求的内容w, err := c.readRequest()
- 判断handler是否为空，如果没有设置handler（这个例子就没有设置handler），handler就设置为DefaultServeMux
- 调用handler的ServeHttp
- 根据request选择handler，并且进入到这个handler的ServeHTTP
- 选择handler：
	* 判断是否有路由能满足这个request（循环遍历ServerMux的muxEntry）
	* 如果有路由满足，调用这个路由handler的ServeHttp
	* 如果没有路由满足，调用NotFoundHandler的ServeHttp

## 其它特性说明
本次作业在原有的基础上优化了登录页的页面布局以及实现方式，同时，增加了一个index页面和error页的跳转，实现了识别用户名和密码，根据用户名和密码跳转不同页面。