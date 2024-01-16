package createidx

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"txtresearch/gstr"
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
	filetext, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	loop := 0
	filename := gstr.Do(path, "\\", ".txt", true, false)
	filename = gstr.RStr(filename, "-")
	filename = "《" + filename + "》\n" //加入经名并加上"《"，以便专搜经名
	secs := d.splitText(filename + string(filetext))
	lenstr := strconv.Itoa(len(secs))
	mlenstr := strconv.Itoa(len(lenstr))

	var secNo string
	for _, sv := range secs {
		if sv == "" { //根据这个\n，可以判断文章的段落
			sv = "\n"
		}
		loop++
		secNo = fmt.Sprintf("%0"+mlenstr+"d", loop)
		d.table["c"].Ins( //添加到表c
			map[string]string{
				"id": id + "," + secNo,
				"s":  sv,
				"r":  "", // pubgo.Reverse(sv),
			})
	}
}

// 按。将文章分段
func (d *DzjFile) splitText(filetext string) []string {
	filetext = strings.Replace(filetext, "\r", "\r\n", -1)
	filetext = strings.Replace(filetext, "\t", "\t\n", -1)
	filetext = strings.Replace(filetext, "。", "。\n", -1)
	filetext = strings.Replace(filetext, ".", ".\n", -1)
	filetext = strings.Replace(filetext, "？", "？\n", -1)
	filetext = strings.Replace(filetext, "！", "！\n", -1)
	filetext = strings.Replace(filetext, "?", "?\n", -1)
	filetext = strings.Replace(filetext, "!", "!\n", -1)
	filetext = strings.Replace(filetext, "；", "；\n", -1)
	filetext = strings.Replace(filetext, ";", ";\n", -1)
	return strings.Split(filetext, "\n")
}
