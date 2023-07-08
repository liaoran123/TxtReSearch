package routers

import (
	"fmt"
	"txtresearch/xbdb"
)

var Table map[string]*xbdb.Table
var xb *xbdb.Xb
var loop int

func Ini() {
	//打开或创建数据库
	dbpath := ConfigMap["dbpath"].(string)
	xb = xbdb.NewDb(dbpath + "db")

	//建表
	dbinfo := xbdb.NewTableInfoNil(xb.Db)
	//dbinfo.Del("ca")
	//dbinfo.Del("art")
	//dbinfo.Del("c")

	if dbinfo.GetInfo("c").FieldType == nil {
		createc(dbinfo)
	}

	//打开表操作结构

	/*
		fmt.Println("开始", time.Now())
		copytb()
		fmt.Println("结束", time.Now())
	*/
	Table = xb.GetTables()

	//Table["c"].Select.ForDbase(Pr)

	//Table["ca"].Select.FindPrefixFun([]byte("ca."), false, Pr) //ca.fid-\x00\x00\x00\x01-
	/*
		for i := 10; i < 49; i++ {
			id := strconv.Itoa(i)
			xbdb.Tables["record"].Del(id)
		}*/
	//目录入加载内存
	//CRAMs = NewCataRAMs()
	//CRAMs.LoadCataRAM()
	//文章对应的目录fid加载入内存
	//LoadartRAM()
}
func Pr(k, v []byte) bool {
	ks, kv := string(k), string(v)
	fmt.Println(ks, kv)
	return true
}

// 创建文章内容表，该表是全文搜索，故而名称尽量短，可以减少文件大小。
// 带全文搜索索引的内容表c
func createc(tbifo *xbdb.TableInfo) {
	name := "c"                                         //目录表，
	fields := []string{"id", "s", "r"}                  //字段 s 是文章的分段内容，r，是s的反转字符串，用于前置匹配词,pos,为位置
	fieldType := []string{"string", "string", "string"} //字段
	idxs := []string{}                                  //索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔
	fullText := []string{"1", "2"}                      //考据级全文搜索索引字段的下标。
	ftlen := "7"                                        //全文搜索的长度，中文默认是7
	patterns := []string{"1", "5"}                      //搜索词模型。 1,中文;2字母;3，数字；4，标点符号；5，自定义。不符合的字被过滤。可以组合。
	diychar := "《》"
	r := tbifo.Create(name, ftlen, diychar, fields, fieldType, idxs, fullText, patterns)
	fmt.Printf("r: %v\n", r)
}
