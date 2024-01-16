package main

import (
	"fmt"
	"topsdk"
	"topsdk/ability370"
	"topsdk/ability370/request"
)

func mai1n() {
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
	req.SetStartDsr(10)
	req.SetPageSize(20)
	req.SetPageNo(1)
	req.SetPlatform(1)
	//req.SetEndTkRate(1234) //商品筛选-淘客佣金比率上限。如：1234表示12.34%
	//req.SetStartTkRate(1234) //商品筛选-淘客佣金比率下限。如：1234表示12.34%
	//req.SetEndPrice(10) //商品筛选-折扣价范围上限。单位：元
	//req.SetStartPrice(10) //商品筛选-折扣价范围下限。单位：元
	req.SetIsOverseas(false)   //商品筛选-是否海外商品。true表示属于海外商品，false或不设置表示不限
	req.SetIsTmall(false)      //商品筛选-是否天猫商品。true表示属于天猫商品，false或不设置表示不限
	req.SetSort("tk_rate_des") //排序_des（降序），排序_asc（升序），销量（total_sales），淘客佣金比率（tk_rate）， 累计推广量（tk_total_sales），总支出佣金（tk_total_commi），价格（price），匹配分（match）
	//req.SetItemloc("杭州") //商品筛选-所在地
	//req.SetCat("16,18") //商品筛选-后台类目ID。用,分割，最大10个，该ID可以通过taobao.itemcats.get接口获取到
	req.SetQ("女装")
	req.SetMaterialId(2836)       //不传时默认物料id=2836；如果直接对消费者投放，可使用官方个性化算法优化的搜索物料id=17004
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
		fmt.Println(err)
	} else {
		fmt.Println(resp.Body)
	}
}
