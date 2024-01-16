package web

import (
	"context"
	"net/http"
)

type pubpar struct {
	Host   string
	Ucook  Ucook
	Ift    bool      //是否转繁体
	TopDir []dirinfo //搜索项目，根据文件夹得出
}

//var pp pubpar

func newpubpar(req *http.Request) pubpar {
	ckft, _ := req.Cookie("ft")
	Ift := false
	if ckft != nil {
		Ift = ckft.Value == "1"
	}
	return pubpar{
		Host:   ConfigMap["http"].(string) + req.Host,
		Ift:    Ift,
		Ucook:  GetCookie(req),
		TopDir: TopDir,
	}
}

func Pubpars() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			ckft, _ := req.Cookie("ft")
			Ift := false
			if ckft != nil {
				Ift = ckft.Value == "1"
			}
			pp := pubpar{
				Host:   ConfigMap["http"].(string) + req.Host,
				Ift:    Ift,
				Ucook:  GetCookie(req),
				TopDir: TopDir,
			}
			data := map[string]interface{}{}
			data["pubpar"] = pp
			ctx := context.WithValue(req.Context(), "data", data) //"上下文数据"
			f(w, req.WithContext(ctx))
		}
	}
}
