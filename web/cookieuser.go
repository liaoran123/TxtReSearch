package web

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"
)

const usercookiename = "user"

// 登录用户cookie
type Ucook struct {
	Id, Name string
}

// 设置Cookie
func SetCookie(w http.ResponseWriter, id, name string) {
	//出于安全考虑，每个月换一个cooke名称。
	cookename := usercookiename + time.Now().Month().String()
	value := []byte(id + "," + name)
	cookie := http.Cookie{
		Name:     cookename,
		Value:    base64.URLEncoding.EncodeToString(value),
		Path:     "/",
		Secure:   true, //安全机制，意味着保持Cookie通信只限于加密传输，指示浏览器仅仅在通过安全/加密连接才能使用该Cookie。
		HttpOnly: true, //指示浏览器不要在除HTTP（和 HTTPS)请求之外暴露Cookie。一个有HttpOnly属性的Cookie，不能通过非HTTP方式来访问，例如通过调用JavaScript(例如，引用 document.cookie）
		//MaxAge:   31 * 24 * 60 * 60, //31天. 单位：秒 .cookie过期时间，使用相对时间，比如300s
		Expires: time.Now().AddDate(0, 0, 31), // Cookie 有效期设置为31天.cookie过期时间，使用绝对时间。比如2018/10/10 10:10:10
	}
	http.SetCookie(w, &cookie)
}

// 读取Cookie
func ReadCookie(r *http.Request) string {
	//每个星期换一个cooke名称
	cookename := usercookiename + time.Now().Month().String()
	cookie, err := r.Cookie(cookename)
	if err == nil {
		value, _ := base64.URLEncoding.DecodeString(cookie.Value)
		return string(value)
	}
	return ""
}
func GetCookie(r *http.Request) (uk Ucook) {
	value := ReadCookie(r)
	if value != "" {
		//uk:=ucook{}
		Cookies := strings.Split(value, ",")
		uk.Id = Cookies[0]
		uk.Name = Cookies[1]
	}
	return
}

// 删除Cookie
func DelCookie(w http.ResponseWriter, r *http.Request) {
	cookename := usercookiename + time.Now().Month().String()
	//value := []byte("")
	cookie := http.Cookie{
		Name:  cookename,
		Value: "", //base64.URLEncoding.EncodeToString(value),
		Path:  "/",
		//Secure:   true,                         //安全机制，意味着保持Cookie通信只限于加密传输，指示浏览器仅仅在通过安全/加密连接才能使用该Cookie。
		HttpOnly: true,                         //指示浏览器不要在除HTTP（和 HTTPS)请求之外暴露Cookie。一个有HttpOnly属性的Cookie，不能通过非HTTP方式来访问，例如通过调用JavaScript(例如，引用 document.cookie）
		MaxAge:   -1,                           //31天. 单位：秒 .cookie过期时间，使用相对时间，比如300s
		Expires:  time.Now().AddDate(0, 0, -1), // Cookie 有效期设置为31天.cookie过期时间，使用绝对时间。比如2018/10/10 10:10:10
	}
	http.SetCookie(w, &cookie)
	w.Write([]byte("OK"))
}

/*
cookie.Path("/WEB16");
    代表访问WEB16应用中的任何资源都携带cookie
cookie.Path("/WEB16/cookietest");
    代表访问WEB16中的cookietest时才携带cookie信息
cookie.Domain(".foo.com");
    这对foo.com域下的所有主机都生效(如www.foo.com)，但不包括子域www.abc.foo.com
*/
