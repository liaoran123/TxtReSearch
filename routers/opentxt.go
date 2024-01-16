package routers

import (
	"fmt"
	"net/http"
	"os"
	"txtresearch/createidx"
)

// 打开txt文件
func OpenTxtApi(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	id := req.URL.Query().Get("id")
	fpath := createidx.Dirs[id] + ".txt"
	filetext, err := os.ReadFile(fpath)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(filetext)
}
