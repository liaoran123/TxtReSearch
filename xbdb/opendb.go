// 小白数据库
// 表信息
package xbdb

import (
	"fmt"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var xb *leveldb.DB

// 创建或打开数据库
func OpenDb(fp string) error {
	xb, err = leveldb.OpenFile(fp+"db", nil)
	return err
}

// 创建所有表的操作结构
func OpenTableStructs() map[string]*Table {
	Tables := make(map[string]*Table)
	tbnames := GetTbnames()
	for _, v := range tbnames {
		Tables[v] = NewTable(xb, v)
	}
	return Tables
	/*
		iter := xb.NewIterator(util.BytesPrefix([]byte(Tbspfx+Split)), nil)
		Tables := make(map[string]*Table)
		tbname := ""
		for iter.Next() {
			tbname = strings.Split(string(iter.Key()), Split)[1]
			Tables[tbname] = NewTable(tbname)
		}
		iter.Release()
		if iter.Error() != nil {
			fmt.Printf("iter.Error(): %v\n", iter.Error())
		}*/

}

// 获取数据库所有的表名称
func GetTbnames() (r []string) {
	iter := xb.NewIterator(util.BytesPrefix([]byte(Tbspfx+Split)), nil)
	tbname := ""
	for iter.Next() {
		tbname = strings.Split(string(iter.Key()), Split)[1]
		r = append(r, tbname)
	}
	iter.Release()
	if iter.Error() != nil {
		fmt.Printf("iter.Error(): %v\n", iter.Error())
	}
	return r
}
