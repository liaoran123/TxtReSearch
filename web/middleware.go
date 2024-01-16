package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
)

/*
中间件（通常）是一小段代码，它们接受一个请求，对其进行处理，每个中间件只处理一件事情，
完成后将其传递给另一个中间件或最终处理程序，这样就做到了程序的解耦。
//解耦合，其实就是等于可以平衡执行，类似多线程执行，可以互不干涉。
//如路由，相同的入参，执行不同的事情。并且无数据交涉，串联。
//也即是能够做到互不干涉，就做。就是解耦合。
如果没有中间件那么我们必须在最终的处理程序中来完成这些处理操作，这无疑会造成处理程序的臃肿和代码复用率不高的问题。！！！！！！！
中间件的一些常见用例是请求日志记录，Header操纵、HTTP请求认证和ResponseWriter劫持等等。
*/

/*
中间件只将http.HandlerFunc作为其参数，在中间件里将其包装并返回新的http.HandlerFunc供服务器服务复用器调用。
这里我们创建一个新的类型Middleware，这会让最后一起链式调用多个中间件变的更简单。
*/
//创建中间件
type Middleware func(http.HandlerFunc) http.HandlerFunc

// 把应用到http.HandlerFunc处理器的中间件
// 按照先后顺序和处理器本身链起来供http.HandleFunc调用
// --------------为f层递叠加装饰器---------
func DoMiddleware(f http.HandlerFunc, md ...Middleware) http.HandlerFunc {
	for _, hf := range md {
		f = hf(f)
	}
	return f
}

// -------中间件模板---------------
func ModalMiddleware() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			data := req.Context().Value("data") //读取上下文内容
			fmt.Printf("data: %v\n", data)
			if !strings.Contains(req.URL.Path, ".txt") {
				f, _ := os.ReadFile(cdir + req.URL.Path)
				w.Write(f)
				return //如果这里退出，则断开DoMiddleware之后的链路。
			}
			//--------------上面代码是装饰器------------------------
			//--------------上面代码是本身-----------------------
			ctx := context.WithValue(req.Context(), "data", "上下文数据") //"上下文数据"
			f(w, req.WithContext(ctx))
		}
	}
}
