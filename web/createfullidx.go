package web

import (
	"net/http"
)

// 创建大藏经全文索引
func Createfullidx(w http.ResponseWriter, req *http.Request) {
	psw := req.URL.Query().Get("psw")
	if psw == "" { //不设置密码则不能进行操作
		w.Write([]byte("管理员已关闭该权限。"))
		return
	}
	if psw == ConfigMap["psw"].(string) { //管理员密码
		/*
			//请先将原数据库改名以作为备份，下面的功能将重新以db为名创建索引数据库。
			createidx.CreatIdx(Table)
		*/
		w.Write([]byte("暂时屏蔽"))
	}

	w.Write([]byte("ok"))
}
