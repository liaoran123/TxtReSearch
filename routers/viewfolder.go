package routers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"txtresearch/createidx"
)

// 浏览大藏经文件夹
func ViewFoldersApi(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	ph := req.URL.Query().Get("ph")
	path := createidx.Root + "\\" + ph //获取绝对路径
	dirlist, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte(err.Error()))
		return
	}
	jsonstr := "{\"fslist\":["
	for _, v := range dirlist {
		jsonstr += "{\"fname\":\"" + v.Name() + "\"},"
	}
	jsonstr = strings.Trim(jsonstr, ",") + "]}"
	w.Write([]byte(jsonstr))
}
