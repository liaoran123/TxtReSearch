package routers

import (
	"html/template"
	"net/http"
)

func AdminHtml(w http.ResponseWriter, req *http.Request) {
	var t *template.Template
	t, _ = template.ParseFiles("admin/index.html") //从文件创建一个模板
	t.Execute(w, nil)
}

func Arthtml(w http.ResponseWriter, req *http.Request) {
	var t *template.Template
	t, _ = template.ParseFiles("admin/art.html") //从文件创建一个模板
	t.Execute(w, nil)
}
func Catahtml(w http.ResponseWriter, req *http.Request) {
	var t *template.Template
	t, _ = template.ParseFiles("admin/cata.html") //从文件创建一个模板
	t.Execute(w, nil)
}
