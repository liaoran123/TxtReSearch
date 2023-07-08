package routers

import (
	"net/http"

	"strconv"
	"txtresearch/gstr"
	"txtresearch/pubgo"
)

func gettongji() string {
	ifo := ""
	sum := 0
	for k, v := range pubgo.Tj.Tjs {
		sum += v.Bws
		ifo += k + ":" + strconv.Itoa(v.Bws) + "\n"
	}
	return ifo + "总计：" + strconv.Itoa(sum)
}

func Test(w http.ResponseWriter, req *http.Request) {
	//统计
	pubgo.Tj.Brows(gstr.Mstr(req.URL.Path, "/", "/"))

	id := req.URL.Query().Get("id")
	if id == "" {
		w.Write([]byte(""))
		return
	}
	rst := gettongji()
	w.Write([]byte(rst))

}
