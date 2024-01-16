package web

import (
	"net/http"

	"strconv"
	"txtresearch/pubgo"
)

// 统计
func gettongji() string {
	ifo := ""
	sum := 0
	for k, v := range pubgo.Tj.Tjs {
		sum += v.Bws
		ifo += k + ":" + strconv.Itoa(v.Bws) + "\n"
	}
	return ifo + "总计：" + strconv.Itoa(sum)
}

func Tongji(w http.ResponseWriter, req *http.Request) {
	//统计
	pubgo.Tj.Brows("tongji")
	rst := ""
	if req.URL.Query().Get("id") == "" {
		rst = gettongji()
	} else {
		rst = fw()
	}
	w.Write([]byte(rst))
}
func fw() string {
	ifo := ""
	sum, ct := 0, 0
	for k := range cips.tjm {
		ct = cips.tjm[k]
		sum += ct
		ifo += k + "," + strconv.Itoa(ct) + "\n"
	}
	return ifo + "总计：" + strconv.Itoa(sum)
}
