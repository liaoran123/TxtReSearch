package createidx

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"txtresearch/gstr"
	"txtresearch/pubgo"
	"txtresearch/xbdb"
)

// 根据路径打开文本，并建立索引
type DzjFile struct {
	table map[string]*xbdb.Table
}

func NewFileText(Table map[string]*xbdb.Table) *DzjFile {
	return &DzjFile{
		table: Table,
	}
}
func (d *DzjFile) CreateIdx(id, path string) {
	filetext, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	loop := 0
	filename := gstr.Do(path, "-", ".txt", true, false)
	filename = "《" + filename + "》\n" //加入经名并加上"《"，以便专搜经名
	secs := d.splitText(filename + string(filetext))
	lenstr := strconv.Itoa(len(secs))
	mlenstr := strconv.Itoa(len(lenstr))

	var secNo string
	newid := ""
	for _, sv := range secs {
		if sv == "" {
			continue
		}
		loop++
		secNo = fmt.Sprintf("%0"+mlenstr+"d", loop)
		newid = strings.Replace(id, "-", "|", -1) //id=1-01-58-15，由于id的分隔符是xbdb的分隔符，为避免转义造成占用空间，这里更改为id=1|01|58|15
		d.table["c"].Ins(                         //添加到表c
			map[string]string{
				"id": newid + "," + secNo,
				"s":  sv,
				"r":  pubgo.Reverse(sv),
			})
	}
}

// 按。将文章分段
func (d *DzjFile) splitText(filetext string) []string {
	filetext = strings.Replace(filetext, "\r", "\n", -1)
	filetext = strings.Replace(filetext, "\t", "\n", -1)
	filetext = strings.Replace(filetext, "。", "\n", -1)
	filetext = strings.Replace(filetext, ".", "\n", -1)
	filetext = strings.Replace(filetext, "？", "\n", -1)
	filetext = strings.Replace(filetext, "！", "\n", -1)
	filetext = strings.Replace(filetext, "?", "\n", -1)
	filetext = strings.Replace(filetext, "!", "\n", -1)
	filetext = strings.Replace(filetext, "；", "\n", -1)
	filetext = strings.Replace(filetext, ";", "\n", -1)
	return strings.Split(filetext, "\n")
}
