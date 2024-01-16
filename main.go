// 大藏经定制版考据级搜索引擎
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"txtresearch/pubgo"
	"txtresearch/web"

	"github.com/kardianos/service"
)

// 读取配置文件装入map
func readconfig() {
	pubgo.Ini() //pubgo.GetCurrentAbPath()
	text, _ := os.ReadFile(pubgo.CurPath + "/web/config.json")
	web.ConfigMap = make(map[string]interface{})
	json.Unmarshal(text, &web.ConfigMap)
}

// 初始化全局数据、变量
func ini() {
	web.Ini()
	web.LoadDirs()               //加载全部目录
	pubgo.Tj = pubgo.Newtongji() //统计
}

// 添加路由
func addrouters() {
	/*
		//routers包为api，暂时不开放
				http.HandleFunc("/admin/", routers.AdminHtml)
				http.HandleFunc("/admin/cata/", routers.Catahtml)
				http.HandleFunc("/admin/search/", routers.Searchhtml)
				http.HandleFunc("/admin/art/", routers.Arthtml)


			http.HandleFunc("/api/cata/", routers.Cata) //目录，get,post,put,delete
			http.HandleFunc("/api/art/", routers.Art)   //文章，get,post,put,delete

			http.HandleFunc("/api/search/", routers.Search) //http.HandleFunc("/api/search/", routers.Search) //搜索

			http.HandleFunc("/api/art/item/", routers.Artitem) //获取目录下的文章列表
			http.HandleFunc("/api/art/meta/", routers.Meta)    //获取文章摘录
			http.HandleFunc("/api/Idxpfx/", routers.Idxpfx)    //搜索词为前缀的相关词

			http.HandleFunc("/test", routers.Test)
	*/
	http.HandleFunc("/static123/", web.Static) //静态文件服务器
	http.HandleFunc("/", web.Temp)
	//http.HandleFunc("/home/", web.Index) //首页
	http.HandleFunc("/home/", web.Index()) //首页
	//不使用微服务模式，是为了方便打包，简化其他人搭建使用
	//http.HandleFunc("/createfullidx/", web.Createfullidx) //创建数据
	//http.HandleFunc("/skqs/", web.ViewFolders)            //浏览文件夹
	http.HandleFunc("/s/", web.Search())
	http.HandleFunc("/art/", web.Art)
	http.HandleFunc("/dir/", web.Dir)
	http.HandleFunc("/getonejuzi/", web.Getonejuzi) //获取文章摘录
	http.HandleFunc("/metas/", web.Metas)
	http.HandleFunc("/meta/", web.Meta)
	http.HandleFunc("/idxpfx/", web.Idxpfx) //搜索词为前缀的相关词
	http.HandleFunc("/kw/", web.Kws)
	http.HandleFunc("/pcv/", web.Pcv)
	http.HandleFunc("/taobao/", web.Taobao)

	http.HandleFunc("/shuoming/", web.Shuoming)
	http.HandleFunc("/zanzhu/", web.Zanzhu)
	http.HandleFunc("/updates/", web.Updates)
	http.HandleFunc("/research/", web.Research())
	http.HandleFunc("/ft/", web.Ft)         //设置繁简体
	http.HandleFunc("/tongji/", web.Tongji) //统计
	http.HandleFunc("/serror/", web.Serror)
	http.HandleFunc("/kwas/", web.Kwas)
	//http.HandleFunc("/mid/", Chain(Hello, Method("GET"), Logging())) //中间件测试
}

// 运行服务
func run() {
	fmt.Println(time.Now())
	fmt.Println("四库全书txtReSearch服务器程序启动成功!")
	fmt.Println("----------------------------")
	port := web.ConfigMap["port"].(string) //从配置文件获取port
	fmt.Println("请在浏览器打开下面的地址即可使用")
	fmt.Println("http://127.0.0.1:" + port)
	fmt.Println("----------------------------")
	//Openurl("http://127.0.0.1:" + port) //打开网址，用于用户本地搭建使用
	err := http.ListenAndServe(":"+port, nil)
	//log.Fatal(err)
	if err != nil {
		fmt.Println("请更正错误后重启程序：", err)
	}
}
func main() {
	//mai1n()
	svcConfig := &service.Config{
		Name:        "TxtReSearch service",
		DisplayName: "TxtReSearch",
		Description: "四库全书考据级搜索服务.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	//命令行参数安装服务
	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			x := s.Install()
			if x != nil {
				fmt.Println("error:", x.Error())
				return
			}
			fmt.Println("服务安装成功")
			return
		} else if os.Args[1] == "uninstall" {
			x := s.Uninstall()
			if x != nil {
				fmt.Println("error:", x.Error())
				return
			}
			fmt.Println("服务卸载成功")
			return
		}
	}
	//s.Uninstall()
	//自动安装服务
	//第一次执行一次。win服务，所有读取文件的路径都需要是绝对路径，否则无法读取。
	fname := pubgo.GetCurrentAbPath() + "/config.txt"
	text, err := os.ReadFile(fname)
	if err != nil {
		println(err.Error())
	}
	fp := string(text)
	cfp, _ := os.Executable()
	if fp == "" { //第一次文件是空，故而安装服务
		if s.Install() == nil {
			os.WriteFile(fname, []byte(cfp), 0666)
			println("安装服务成功")
		}
	} else { //不是第一次，则判断用户是否更改了目录或修改了名称等
		if fp != cfp { //先卸载后安装
			if s.Uninstall() == nil {
				if s.Install() == nil {
					println("重装服务成功")
				}
			}
			os.WriteFile(fname, []byte(cfp), 0666)
		}
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}

}
