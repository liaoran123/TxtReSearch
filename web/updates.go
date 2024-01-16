package web

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"strings"
)

// 用于对网站数据热更新
func Updates(w http.ResponseWriter, req *http.Request) {
	psw := req.URL.Query().Get("psw")
	if psw == "" { //不设置密码则不能进行操作
		w.Write([]byte("ok"))
		return
	}
	cpsw := ConfigMap["psw"].(string)
	if cpsw == "" {
		w.Write([]byte("无权限操作！"))
		return
	}
	if psw == cpsw { //管理员密码
		tempdo()
		println("1")
	} else {
		w.Write([]byte("密码不对！"))
		return
	}

	w.Write([]byte("ok"))
}

func tempdo() {
	var inspar map[string]string
	file, err := os.Open(cdir + "/word/all.txt")
	if err != nil {
		println(err.Error())
		return
	}
	defer file.Close()
	var ps []string
	sline := ""
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			return
		}
		sline = string(line)
		sline = strings.Replace(sline, " ", "", -1)
		if sline == "\t" || sline == "" {
			continue
		}
		ps = strings.Split(sline, "\t")
		inspar = map[string]string{}
		inspar["id"] = ps[0]
		//inspar["key"] = ps[0]
		//inspar["count"] = ps[1]
		r := Table["kw"].Ins(inspar)
		if !r.Succ {
			println(r.Info)
		}
	}
}
