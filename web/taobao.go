package web

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"topsdk"
	"topsdk/ability370"
	"topsdk/ability370/request"
	"txtresearch/gstr"
	"txtresearch/pubgo"
)

type taobao struct {
	Pubpar        pubpar
	Kw            string //搜索词
	Shop          []shop
	Total_results string //搜索到符合条件的结果总数
	P             int
	Pn            int
}
type shop struct {
	Title, Shop_title             string   //标题
	Pict_url                      string   //商品信息-商品主图
	Small_images                  []string //商品小图列表
	Reserve_price, Zk_final_price string   //商品一口价格,折扣价（元） 若属于预售商品，付定金时间内，折扣价=预售价
	Category_id, Category_name    string
	User_type                     string //店铺信息-卖家类型。0表示集市，1表示天猫
	Provcity                      string //商品信息-宝贝所在地
	Nick                          string //店铺信息-卖家昵称
	Coupon_info                   string //优惠券信息-优惠券满减信息
	Coupon_share_url              string //链接-宝贝+券二合一页面链接
	Url                           string //链接-宝贝推广链接
}

func Taobao(w http.ResponseWriter, req *http.Request) {

	rd := new(taobao)
	rd.Pubpar = newpubpar(req)
	rd.Kw = req.URL.Query().Get("kw")
	p := req.URL.Query().Get("p")
	intp, err := strconv.Atoi(p)
	if err != nil {
		intp = 1
	}
	rd.P = intp
	rd.Pn = intp + 1
	ca := req.URL.Query().Get("ca")
	shops := gettaobao(rd.Kw, ca, int64(intp))
	rd.Total_results = gstr.Mstr(shops, "total_results\":", ",")
	sstr := gstr.Do(shops, "[", "]", false, true)
	shopstrlist := strings.Split(sstr, "},")
	for _, v := range shopstrlist {
		ishop := shop{}
		ishop.Title = getshopitemval(v, "title")
		ishop.Shop_title = getshopitemval(v, "shop_title")
		ishop.Pict_url = getshopitemval(v, "pict_url")
		ishop.Pict_url = strings.Replace(ishop.Pict_url, "\\/", "/", -1)
		small_images := gstr.Mstr(v, "\"small_images\":[", "]") //getshopitemval(v, "small_images")
		small_images = strings.Replace(small_images, "\"", "", -1)
		small_images = strings.Replace(small_images, "\\/", "/", -1)
		//small_images = strings.Trim(small_images, "[")
		//small_images = strings.Trim(small_images, "]")
		ishop.Small_images = strings.Split(ishop.Pict_url+","+small_images, ",")
		ishop.Reserve_price = getshopitemval(v, "reserve_price")
		ishop.Zk_final_price = getshopitemval(v+",", "zk_final_price")
		ishop.Zk_final_price = strings.Trim(ishop.Zk_final_price, "}")
		ishop.Category_id = getshopitemval(v+",", "category_id")
		ishop.Category_name = getshopitemval(v+",", "category_name")
		ishop.Category_name = strings.Replace(ishop.Category_name, "\\/", "/", -1)
		ishop.User_type = getshopitemval(v, "user_type")
		ishop.Provcity = getshopitemval(v, "provcity")
		ishop.Nick = getshopitemval(v, "nick")
		ishop.Coupon_info = getshopitemval(v, "coupon_info")
		ishop.Coupon_share_url = getshopitemval(v, "coupon_share_url")
		ishop.Coupon_share_url = strings.Replace(ishop.Coupon_share_url, "\\/", "/", -1)
		ishop.Url = getshopitemval(v, "url")
		ishop.Url = strings.Replace(ishop.Url, "\\/", "/", -1)
		rd.Shop = append(rd.Shop, ishop)
	}

	pubgo.Tj.Brows("taobao")

	//--组织模板数据
	tpl := cdir + "/tpl"
	if ConfigMap["tpl"] != nil {
		tpl = ConfigMap["tpl"].(string)
	}
	TemplatesFiles := []string{
		tpl + "/taobao.html",
		tpl + "/pub/static.html",
		tpl + "/pub/header.html",
		tpl + "/pub/search.html",
		tpl + "/pub/gomove.html",
		tpl + "/pub/footer.html", // 多加的文件
	}

	funcMap := template.FuncMap{ //--需要注册的函数
		"jf":     pubgo.Jf,
		"tohtml": Tohtml,
	}
	t, err := template.New("taobao.html").Funcs(funcMap).ParseFiles(TemplatesFiles...)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	err = t.ExecuteTemplate(w, "taobao.html", rd)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

// 淘客api
func gettaobao(kw, ca string, p int64) string {
	client := topsdk.NewDefaultTopClient("34622800", "aa37cf83a102477d4c1313f30c7dc927", "http://gw.api.taobao.com/router/rest", 20000, 20000)
	ability := ability370.NewAbility370(&client)
	/*
		//物料评估-商品列表
		taobaoTbkDgMaterialOptionalUcrowdrankitems := domain.TaobaoTbkDgMaterialOptionalUcrowdrankitems{}
		taobaoTbkDgMaterialOptionalUcrowdrankitems.SetCommirate(1234)
		taobaoTbkDgMaterialOptionalUcrowdrankitems.SetPrice("10.12")
		taobaoTbkDgMaterialOptionalUcrowdrankitems.SetItemId("542808901898")
	*/
	req := request.TaobaoTbkDgMaterialOptionalRequest{}
	req.SetStartDsr(10) //商品筛选(特定媒体支持)-店铺dsr评分。筛选大于等于当前设置的店铺dsr评分的商品0-50000之间
	req.SetPageSize(20) //页大小，默认20，1~100
	req.SetPageNo(p)    //第几页，默认：１
	req.SetPlatform(1)  //链接形式：1：PC，2：无线，默认：１
	//req.SetEndTkRate(1234) //商品筛选-淘客佣金比率上限。如：1234表示12.34%
	//req.SetStartTkRate(1234) //商品筛选-淘客佣金比率下限。如：1234表示12.34%
	//req.SetEndPrice(10) //商品筛选-折扣价范围上限。单位：元
	//req.SetStartPrice(10) //商品筛选-折扣价范围下限。单位：元
	req.SetIsOverseas(false)       //商品筛选-是否海外商品。true表示属于海外商品，false或不设置表示不限
	req.SetIsTmall(false)          //商品筛选-是否天猫商品。true表示属于天猫商品，false或不设置表示不限
	req.SetSort("total_sales_des") //排序_des（降序），排序_asc（升序），销量（total_sales），淘客佣金比率（tk_rate）， 累计推广量（tk_total_sales），总支出佣金（tk_total_commi），价格（price），匹配分（match）
	//req.SetItemloc("杭州") //商品筛选-所在地
	req.SetCat(ca) //商品筛选-后台类目ID。用,分割，最大10个，该ID可以通过taobao.itemcats.get接口获取到
	req.SetQ(kw)
	req.SetMaterialId(17004)      //不传时默认物料id=2836；如果直接对消费者投放，可使用官方个性化算法优化的搜索物料id=17004
	req.SetHasCoupon(true)        //优惠券筛选-是否有优惠券。true表示该商品有优惠券，false或不设置表示不限
	req.SetIp("13.2.33.4")        //ip参数影响邮费获取，如果不传或者传入不准确，邮费无法精准提供
	req.SetAdzoneId(115524200355) //mm_xxx_xxx_12345678三段式的最后一段数字,推广广告位
	req.SetNeedFreeShipment(true) //商品筛选-是否包邮。true表示包邮，false或不设置表示不限
	req.SetNeedPrepay(true)       //商品筛选-是否加入消费者保障。true表示加入，false或不设置表示不限
	//req.SetIncludePayRate30(true) //商品筛选(特定媒体支持)-成交转化是否高于行业均值。True表示大于等于，false或不设置表示不限
	//req.SetIncludeGoodRate(true)  //商品筛选-好评率是否高于行业均值。True表示大于等于，false或不设置表示不限
	//req.SetIncludeRfdRate(true)   //商品筛选(特定媒体支持)-退款率是否低于行业均值。True表示大于等于，false或不设置表示不限
	req.SetNpxLevel(2)          //商品筛选-牛皮癣程度。取值：1不限，2无，3轻微
	req.SetDeviceEncrypt("MD5") //智能匹配-设备号加密类型：MD5
	req.SetDeviceValue("xxx")
	req.SetDeviceType("IMEI")
	//req.SetEndKaTkRate(1234) //商品筛选-KA媒体淘客佣金比率上限。如：1234表示12.34%
	//req.SetStartKaTkRate(1234) //商品筛选-KA媒体淘客佣金比率下限。如：1234表示12.34%
	//req.SetLockRateEndTime(1567440000000) //锁佣结束时间
	//req.SetLockRateStartTime(1567440000000)//锁佣开始时间
	//req.SetLongitude("121.473701")
	//req.SetLatitude("31.230370")
	//req.SetCityCode("310000")
	//req.SetSellerIds("1,2,3,4")
	//req.SetSpecialId("2323")
	//req.SetRelationId("3243")
	//req.SetPageResultKey("abcdef")
	//req.SetUcrowdId(1)

	//req.SetUcrowdRankItems([]domain.TaobaoTbkDgMaterialOptionalUcrowdrankitems{})
	//req.SetGetTopnRate(0)
	//req.SetBizSceneId("1")
	//req.SetPromotionType("2")

	resp, err := ability.TaobaoTbkDgMaterialOptional(&req)
	if err != nil {
		return ""
	} else {
		return resp.Body
	}

}

func getshopitemval(shops, filename string) string {
	val := gstr.Mstr(shops, "\""+filename+"\":", ",")
	val = strings.Replace(val, "\"", "", -1)
	return val
}
func Tohtml(text string) template.HTML {
	return template.HTML(text)
}
