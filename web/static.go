package web

import (
	"net/http"
)

// --静态文件服务
func Static(w http.ResponseWriter, req *http.Request) {
	had := http.StripPrefix("/static123/", http.FileServer(http.Dir(cdir+"/static")))
	had.ServeHTTP(w, req)
}
