package web

import (
	"net/http"
	"os"
	"strings"
	"txtresearch/pubgo"
)

// 由于网站被攻击，临时使用防御
func Temp(w http.ResponseWriter, req *http.Request) {
	if strings.Contains(req.URL.Path, ".txt") {
		f, _ := os.ReadFile(cdir + req.URL.Path)
		w.Write(f)
		return
	}
	pubgo.Tj.Brows("Temp")
	http.ServeFile(w, req, cdir+"/tpl/temp.html")
}
