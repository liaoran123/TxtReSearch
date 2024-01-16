package web

import (
	"encoding/json"
	"net/http"
	"time"
	"txtresearch/xbdb"
)

// 设置繁简体
func Ft(w http.ResponseWriter, req *http.Request) {
	var r xbdb.ReInfo
	ft := req.URL.Query().Get("l")
	cookie := http.Cookie{
		Name:    "ft",
		Value:   ft,
		Path:    "/",
		Expires: time.Now().AddDate(0, 0, 31), // Cookie 有效期设置为31天.cookie过期时间，使用绝对时间。比如2018/10/10 10:10:10
	}
	http.SetCookie(w, &cookie)
	r.Succ = true
	r.Info = "ok"
	json.NewEncoder(w).Encode(r)
}
