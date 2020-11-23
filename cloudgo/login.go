package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/unrolled/render"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析url传递的参数，对于POST则解析响应包的主体（request body）
	//注意:如果没有调用ParseForm方法，下面无法获取表单的数据
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}

// func login(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("method:", r.Method) //获取请求的方法
// 	if r.Method == "GET" {
// 		t, _ := template.ParseFiles("login.html")
// 		log.Println(t.Execute(w, nil))
// 	} else {
// 		//请求的是登录数据，那么执行登录的逻辑判断
// 		fmt.Println("username:", r.Form["username"])
// 		fmt.Println("password:", r.Form["password"])
// 	}
// }
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
