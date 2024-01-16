package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

/*
中间件（通常）是一小段代码，它们接受一个请求，对其进行处理，每个中间件只处理一件事情，
完成后将其传递给另一个中间件或最终处理程序，这样就做到了程序的解耦。
如果没有中间件那么我们必须在最终的处理程序中来完成这些处理操作，这无疑会造成处理程序的臃肿和代码复用率不高的问题。！！！！！！！
中间件的一些常见用例是请求日志记录，Header操纵、HTTP请求认证和ResponseWriter劫持等等。
*/

/*
中间件只将http.HandlerFunc作为其参数，在中间件里将其包装并返回新的http.HandlerFunc供服务器服务复用器调用。
这里我们创建一个新的类型Middleware，这会让最后一起链式调用多个中间件变的更简单。
*/
//创建中间件
type Middleware func(http.HandlerFunc) http.HandlerFunc

/*
中间件是使用装饰器模式实现的，下面的中间件通用代码模板让我们平时编写中间件变得更容易，我们在自己写中间件的时候只需要往样板里填充需要的代码逻辑即可。
*/
//中间件代码模板
func CreateNewMiddleware() Middleware {
	// 创建一个新的中间件
	middleware := func(next http.HandlerFunc) http.HandlerFunc {
		// 创建一个新的handler包裹next
		handler := func(w http.ResponseWriter, r *http.Request) {

			// 中间件的处理逻辑
			// 调用下一个中间件或者最终的handler处理程序
			next(w, r)
		}

		// 返回新建的包装handler
		return handler
	}

	// 返回新建的中间件
	return middleware
}

// 使用中间件
// 记录每个URL请求的执行时长
func Logging() Middleware {

	// 创建中间件
	return func(f http.HandlerFunc) http.HandlerFunc {

		// 创建一个新的handler包装http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// 中间件的处理逻辑
			start := time.Now()

			ctx := context.WithValue(r.Context(), "data", "使用上下文传递数据1")
			defer func() { log.Println(r.URL.Path, time.Since(start)) }()
			// 调用下一个中间件或者最终的handler处理程序
			f(w, r.WithContext(ctx))
		}
	}
}

// 验证请求用的是否是指定的HTTP Method，不是则返回 400 Bad Request
func Method(m string) Middleware {

	return func(f http.HandlerFunc) http.HandlerFunc {

		return func(w http.ResponseWriter, r *http.Request) {

			if r.Method != m {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			ctx := context.WithValue(r.Context(), "data", "使用上下文传递数据2")
			f(w, r.WithContext(ctx))
			log.Println(r.URL.Path, "Method")
		}
	}
}

// 把应用到http.HandlerFunc处理器的中间件
// 按照先后顺序和处理器本身链起来供http.HandleFunc调用
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

// 最终的处理请求的http.HandlerFunc
func Hello(w http.ResponseWriter, r *http.Request) {
	data := r.Context().Value("data") //读取上下文内容
	fmt.Fprintln(w, data)
}

func Mainmmidd() {
	//总结，就是解耦合的方法，一个HandlerFunc做一件事情，然后串联起来执行。
	//虽然如此，Chain函数还是觉得有点奇怪。
	http.HandleFunc("/", Chain(Hello, Method("GET"), Logging()))
	http.ListenAndServe(":8080", nil)
}

//解耦合，其实就是等于可以平衡执行，类似多线程执行，可以互不干涉。
//也即是能够做到互不干涉，就做。就是解耦合。
//上面的解耦合，是相同的输入，不同的处理执行，并且互不干涉,干扰。

// ------------------范例2-----------------------------------------
// ContextMiddleware 传递公共参数中间件
func ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "data", "上下文数据") //"上下文数据"，生产中则是返回页面的结构
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
