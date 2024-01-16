package web

import (
	"bytes"
	"net/http"
)

func Getonejuzi(w http.ResponseWriter, req *http.Request) {
	//w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	//w.Header().Set("Content-Type", "application/json")
	const (
		tbname   = "c"
		idxfield = "id"
	)
	params := getparas(req)
	id := params["artid"] + "," + params["secid"]
	bid := []byte(id) //xbdb.SplitToCh([]byte(id))
	key := Table[tbname].Select.GetPkKey(bid)
	text := Table[tbname].Select.GetValue(key)
	//text = bytes.Replace(text, []byte("_"), []byte(""), -1)
	text = bytes.Trim(text, "_")
	w.Write(text)
}
