package global

import "ai-software-copyright-server/internal/application/model/common"

var DefaultWelfate = []common.WelfareLabel{
	{
		Label: "美团外卖",
		Items: []common.WelfareItem{
			{
				Title:     "[外卖]每日限时大额红包",
				Desc:      "外卖红包天天领，拼手气开惊喜红包。",
				ImageUrl:  "https://s3plus.sankuai.com/v1/mss_5017c592a8a946d2a54eb62a76ba299c/nebulafile/910fa09a310aadd229e90e4ad872d86e.png",
				TargetUrl: "http://dpurl.cn/4GIufMiz",
				WxAppid:   "wxde8ac0a21135c07d",
				WxPath:    "/index/pages/h5/h5?weburl=https%3A%2F%2Fclick.meituan.com%2Ft%3Ft%3D1%26c%3D2%26p%3DNWqXdL9zsQgZ",
			}, {
				Title:     "[外卖]美团霸王餐",
				Desc:      "下单试吃得返现，最低1元吃外卖。",
				ImageUrl:  "https://p0.meituan.net/bizad/425eb22a3c3503d04e960506ad618643.jpg",
				TargetUrl: "https://kzurl20.cn/7iRJx6",
				WxAppid:   "wxde8ac0a21135c07d",
				WxPath:    "/waimai/pages/web-view/web-view?type=REDIRECT&webviewUrl=https%3A%2F%2Foffsiteact.meituan.com%2Fact%2Fcps%2Fpromotion%3Fp%3D9bc7499724d849eeb8a3832e11af8373&utm_content=0c3bfd35279b4140b3bd8ecbc41301d6__ac118d53892c424cbd2b0689f95bafd6",
			}, {
				Title:     "[果蔬]商超果蔬红包天天领",
				Desc:      "每周三在线满43减13元大额神券。",
				ImageUrl:  "https://p0.meituan.net/marketingcpsmedia/13aa8ec66decc6c686d9b91ce28659bd143877.jpg",
				TargetUrl: "http://dpurl.cn/ncTMNrJz",
				WxAppid:   "wxde8ac0a21135c07d",
				WxPath:    "/index/pages/h5/h5?weburl=https%3A%2F%2Fclick.meituan.com%2Ft%3Ft%3D1%26c%3D2%26p%3Dbn3NQr9zOozz",
			},
		},
	}, {
		Label: "饿了么",
		Items: []common.WelfareItem{
			{
				Title:    "[外卖]每日限时大额红包",
				Desc:     "热热热，外卖每日限时快闪红包。",
				ImageUrl: "https://img.alicdn.com/imgextra/i2/6000000002139/O1CN01kbFO5p1RfiArcapgD_!!6000000002139-2-o2oad.png",
				WxAppid:  "wxece3a9a4c82f58c9",
				WxPath:   "ele-recommend-price/pages/guest/index?scene=7df77ebc0461418abc30adc312ba86a9",
			}, {
				Title:     "[外卖]天天领红包",
				Desc:      "外卖红包天天领，最高抢66元大额红包。",
				ImageUrl:  "https://img.alicdn.com/imgextra/i2/6000000005066/O1CN01ky8Jg81nIHHNDoI7P_!!6000000005066-2-o2oad.png",
				TargetUrl: "https://u.ele.me/ZV5JqjZA",
				WxAppid:   "wxece3a9a4c82f58c9",
				WxPath:    "commercialize/pages/taoke-guide/index?scene=068b5cf8dda4434e87abc838729b3d16",
			}, {
				Title:     "[外卖]消费日外卖专享红包",
				Desc:      "城市消费日，大额红包天天领。",
				ImageUrl:  "https://img.alicdn.com/imgextra/i2/6000000004543/O1CN01OrXHFA1jQkB0fb3jG_!!6000000004543-0-o2oad.jpg",
				TargetUrl: "https://u.ele.me/HLLBbbSp",
				WxAppid:   "wxece3a9a4c82f58c9",
				WxPath:    "ad-bdlm-sub/pages/wh-coupon-guide/index?scene=7859ddf930a8484f9d5ea32a85957aae",
			}, {
				Title:    "[果蔬]果蔬零售会场红包",
				Desc:     "买菜买水果，无门槛红包，超低价爆品来袭。",
				ImageUrl: "https://img.alicdn.com/imgextra/i2/6000000003629/O1CN01eZ1REo1cg8EILHbEL_!!6000000003629-0-o2oad.jpg",
				WxAppid:  "wxece3a9a4c82f58c9",
				WxPath:   "commercialize/pages/bdlm-ls-guide/index?scene=b8fe731b004c43989d4eca1e040ce227",
			}, {
				Title:    "[品牌]天天抢大牌可叠加红包",
				Desc:     "饿了么品牌馆，肯德基、卡斯丁等品牌限时活动，抢5-10元可叠加红包。",
				ImageUrl: "https://img.alicdn.com/imgextra/i4/6000000007953/O1CN01Wlv82i28cWlZF2T3D_!!6000000007953-0-o2oad.jpg",
				WxAppid:  "wxece3a9a4c82f58c9",
				WxPath:   "ad-bdlm-sub/pages/brands-activity/index?scene=dbc26938930b4b13a2aa2f3a21aa969d",
			}, {
				Title:    "[外卖]超抢手特惠会场",
				Desc:     "爆款好价天天有。",
				ImageUrl: "https://img.alicdn.com/imgextra/i4/6000000000383/O1CN01KwzO0c1EhSiU2bsg5_!!6000000000383-0-o2oad.jpg",
				WxAppid:  "wxece3a9a4c82f58c9",
				WxPath:   "ad-bdlm-sub/pages/daily-special-price-foods-guide/index?scene=0cab052e7e604632bf929ccf26338a84",
			}, {
				Title:    "[日用品]大牌爆款特价专享红包",
				Desc:     "大牌爆款，补贴特价，领超值购专享红包。",
				ImageUrl: "https://img.alicdn.com/imgextra/i4/6000000007903/O1CN01WHtCde28Fd6R5QTlD_!!6000000007903-0-o2oad.jpg",
				WxAppid:  "wxece3a9a4c82f58c9",
				WxPath:   "ad-bdlm-sub/pages/daily-special-price-foods-guide/index?scene=1cfdb7c75a18487db2156616dea1dfed",
			}, {
				Title:    "[囤券]囤券券真省钱",
				Desc:     "买东西先囤券，特价好券过期自动退，随时退。",
				ImageUrl: "https://img.alicdn.com/imgextra/i4/6000000000325/O1CN01ZlAzKr1EGtZsHzznp_!!6000000000325-2-o2oad.png",
				WxAppid:  "wxece3a9a4c82f58c9",
				WxPath:   "ad-bdlm-sub/pages/coupon-hoard-guide/index?scene=1fed94db28bc46ebbb037145914f53d7",
			},
		},
	}, {
		Label: "连锁餐饮",
		Items: []common.WelfareItem{
			{
				Title:     "肯德基在线点餐",
				Desc:      "肯德基单单享低价，优惠点餐5折起。",
				ImageUrl:  "https://www.jutuike.com/static/images/kfc.png",
				TargetUrl: "https://www.qipiao.net/h10/#/pages/kfc/kfc?entid=101271&source=APP101271&entpara=266983jutuike123456",
				WxAppid:   "wx89752980e795bfde",
				WxPath:    "/pages/index/index?pub_id=266983&sid=123456&act_id=38&source=jutuike",
			}, {
				Title:     "必胜客线点餐",
				Desc:      "必胜客在线点餐点餐最低7折。",
				ImageUrl:  "https://img.jutuike.com/taokeout/banner/pizzahut_banner.png",
				TargetUrl: "https://kurl02.cn/7NVXTl",
				WxAppid:   "wx89752980e795bfde",
				WxPath:    "/pages/index/index?pub_id=266983&sid=123456&act_id=64&source=jutuike",
			}, {
				Title:     "瑞幸咖啡在线点餐",
				Desc:      "瑞幸咖啡在线点餐，全场饮品5.5折起。",
				ImageUrl:  "https://img.jutuike.com/taokeout/img/luckin.png",
				TargetUrl: "https://5kma.cn/7NVTXm",
				WxAppid:   "wx89752980e795bfde",
				WxPath:    "/pages/index/index?pub_id=266983&sid=123456&act_id=33&source=jutuike",
			}, {
				Title:     "星巴克在线点餐",
				Desc:      "星巴克在线点餐，饮品优惠点餐低至8折起。",
				ImageUrl:  "https://img.jutuike.com/taokeout/img/spk.png",
				TargetUrl: "https://kurl04.cn/7NVTQF",
				WxAppid:   "wx89752980e795bfde",
				WxPath:    "/pages/index/index?pub_id=266983&sid=123456&act_id=34&source=jutuike",
			}, {
				Title:     "奈雪的茶在线点餐",
				Desc:      "奈雪的茶在线点餐，全场优惠8.8折起。",
				ImageUrl:  "https://img.jutuike.com/taokeout/img/nayuki.png",
				TargetUrl: "https://kzurl19.cn/7NVTAd",
				WxAppid:   "wx89752980e795bfde",
				WxPath:    "/pages/index/index?pub_id=266983&sid=123456&act_id=32&source=jutuike",
			}, {
				Title:     "喜茶在线点餐",
				Desc:      "喜茶天天省，全场9.5折。",
				ImageUrl:  "https://img.jutuike.com/taokeout/img/heytea.png",
				TargetUrl: "https://4kma.cn/7NVgzC",
				WxAppid:   "wx89752980e795bfde",
				WxPath:    "/pages/index/index?pub_id=266983&sid=123456&act_id=37&source=jutuike",
			}, {
				Title:     "华莱士在线点餐",
				Desc:      "华莱士在线点餐，低至6折起。",
				ImageUrl:  "https://img.jutuike.com/taokeout/banner/tqwallace_banner.png",
				TargetUrl: "https://kzurl05.cn/7NVSkW",
			}, {
				Title:     "汉堡王在线点餐",
				Desc:      "汉堡王在线点餐，全场8.8折起。",
				ImageUrl:  "https://img.jutuike.com/taokeout/img/burgerking.png",
				TargetUrl: "https://kzurl07.cn/7NVSO0",
				WxAppid:   "wx89752980e795bfde",
				WxPath:    "/pages/index/index?pub_id=266983&sid=123456&act_id=46&source=jutuike",
			}, {
				Title:     "百果园水果外送",
				Desc:      "喜茶天天省，全场9.5折。",
				ImageUrl:  "https://img.jutuike.com/taokeout/img/pagoda.png",
				TargetUrl: "https://kurl06.cn/7NVSeb",
				WxAppid:   "wx89752980e795bfde",
				WxPath:    "/pages/index/index?pub_id=266983&sid=123456&act_id=31&source=jutuike",
			},
		},
	}, {
		Label: "打车出行",
		Items: []common.WelfareItem{
			{
				Title:     "滴滴出行打车劵",
				Desc:      "8折打车券，单笔交易最高抵扣10元，每天都能领。",
				ImageUrl:  "https://img.jutuike.com/taokeout/banner/didi_banner_20240815.jpg",
				TargetUrl: "https://kzurl20.cn/7iRysm",
				WxAppid:   "wxaf35009675aa0b2a",
				WxPath:    "/pages/index/index?scene=rYe2XGg&source_id=266983jutuike123456&ref_from=dunion",
			}, {
				Title:     "花小猪打车劵",
				Desc:      "打车券天天领，最高可领100元券包。",
				ImageUrl:  "https://img.jutuike.com/taokeout/banner/jtk_hxz_banner.png",
				TargetUrl: "https://x.huaxz.cn/pn6WJ5b?source_id=266983jutuike123456",
				WxAppid:   "wxd98a20e429ce834b",
				WxPath:    "/pages/chitu/index?scene=KlqL76g&source_id=266983jutuike123456&ref_from=dunion",
			}, {
				Title:     "T3出行打车劵",
				Desc:      "T3出行超值礼包，最高可领100元券包。",
				ImageUrl:  "https://img.jutuike.com/taokeout/banner/t3go_banner.png",
				TargetUrl: "https://s.t3go.cn/1HnDGjEgW?sourceId=266983jutuike123456",
				WxAppid:   "wxe241a1d8464bc578",
				WxPath:    "independentPackages/webEntry/index?scene=s.t3go.cn/1Q1m4DqYR&sourceId=266983jutuike123456",
			}, {
				Title:     "同程打车打车劵",
				Desc:      "23元打车券免费送，领券打车更优惠。",
				ImageUrl:  "https://img.jutuike.com/taokeout/banner/tcyl_banner_jtk.png",
				TargetUrl: "https://wx.17u.cn/ycoperation/slCoupons?mark=t61198793815p",
				WxAppid:   "wx0d619ddf7b26e48d",
				WxPath:    "/pages/webView/webView?src=https%3A%2F%2Fwx.17u.cn%2Fycoperation%2FslCoupons%3Fmark%3Dt61198793815p",
			}, {
				Title:     "飞猪租车券",
				Desc:      "租车上飞猪，租车最高减免440元。",
				ImageUrl:  "https://img.jutuike.com/taokeout/banner/feizhu_zuche_banner.png",
				TargetUrl: "https://kq-m.dtsoft.cn/#/pages/toMiniProgram/toMiniProgram?&act_id=145&sid=123456&code=sjEaHjBk",
				WxAppid:   "wx6a96c49f29850eb5",
				WxPath:    "pages/home/index?currentTab=rentCar&ttid=12wechat000008258&fpid=18006&fpsid=266983jutuike123456&fp_scene=multi-industry-wx",
			}, {
				Title:     "飞猪酒店机票火车票门票优惠",
				Desc:      "机票酒店火车票门票轻松预订，出行优惠天天享。",
				ImageUrl:  "https://img.jutuike.com/taokeout/banner/feizhu_cps_banner.png",
				TargetUrl: "https://kq-m.dtsoft.cn/#/pages/toMiniProgram/toMiniProgram?&act_id=120&sid=123456&code=sjEaHjBk",
				WxAppid:   "wx6a96c49f29850eb5",
				WxPath:    "pages/home/index?currentTab=hotel&ttid=12wechat000008258&fpid=18006&fpsid=266983jutuike123456&fp_scene=multi-industry-wx",
			},
		},
	}, {
		Label: "电影票",
		Items: []common.WelfareItem{
			{
				Title:    "特价电影票优惠预定",
				Desc:     "电影票在线预订，支持全国连锁影院，最低9.9元购票。",
				ImageUrl: "https://file.youpiaopiao.cn/upload/materials/2022_0414/523434937773326337.jpg",
				WxAppid:  "wx6ce7b07bf7fe6048",
				WxPath:   "/pages/index/index?subplatformid=266983jutuike123456",
			},
		},
	}, {
		Label: "电商购物",
		Items: []common.WelfareItem{
			{
				Title:     "拼多多领券中心",
				Desc:      "拼多多官方领券中心，领券下单更优惠。",
				ImageUrl:  "https://www.jutuike.com/static/images/lingquan.png",
				TargetUrl: "https://p.pinduoduo.com/AoNdKs2M",
				WxAppid:   "wxa918198f16869201",
				WxPath:    "/pages/web/web?specialUrl=1&src=https%3A%2F%2Fmobile.yangkeduo.com%2Fduo_transfer_channel.html%3FresourceType%3D40000%26pid%3D8516041_70976268%26customParameters%3D266983jutuike123456%26authDuoId%3D8516041%26cpsSign%3DCE_241217_8516041_70976268_cb717bda7f217940c50ba3a1c13d62a5%26_x_ddjb_act%3D%257B%2522st%2522%253A%25226%2522%257D%26duoduo_type%3D2",
			}, {
				Title:     "拼多多限时秒杀",
				Desc:      "拼多多超低价好货疯抢。",
				ImageUrl:  "https://www.jutuike.com/static/images/miaosha.png",
				TargetUrl: "https://p.pinduoduo.com/cIaVd4pq",
				WxAppid:   "wxa918198f16869201",
				WxPath:    "/pages/web/web?specialUrl=1&src=https%3A%2F%2Fmobile.yangkeduo.com%2Fduo_transfer_channel.html%3FresourceType%3D4%26pid%3D8516041_141204175%26customParameters%3D266983jutuike123456%26authDuoId%3D8516041%26cpsSign%3DCE_250207_8516041_141204175_c9a3b355715d41dc5ce8053be0de3618%26_x_ddjb_act%3D%257B%2522st%2522%253A%25226%2522%257D%26duoduo_type%3D2",
			}, {
				Title:     "拼多多百亿补贴",
				Desc:      "官方百亿补贴，全场品质保障。",
				ImageUrl:  "https://www.jutuike.com/static/images/butie.png",
				TargetUrl: "https://p.pinduoduo.com/VuhVXeO9",
				WxAppid:   "wxa918198f16869201",
				WxPath:    "/pages/web/web?specialUrl=1&src=https%3A%2F%2Fmobile.yangkeduo.com%2Fduo_transfer_channel.html%3FresourceType%3D39996%26pid%3D8516041_141203721%26_pdd_fs%3D1%26_pdd_tc%3Dffffff%26_pdd_sbs%3D1%26customParameters%3D266983jutuike123456%26authDuoId%3D8516041%26cpsSign%3DCE_241217_8516041_141203721_78f722a1f5aee7abecd473ae8d282ce8%26_x_ddjb_act%3D%257B%2522st%2522%253A%25226%2522%257D%26duoduo_type%3D2",
			},
		},
	}, {
		Label: "特惠酒店",
		Items: []common.WelfareItem{
			{
				Title:     "美团酒店优惠券",
				Desc:      "亲子、学生、商旅等各种活动，领券下单更优惠。",
				ImageUrl:  "https://s3plus.meituan.net/v1/mss_623ff9bb9f0f4bb0a98733e491a41533/union/31abf00bece5497901c87abae391af78.jpg",
				TargetUrl: "https://4kma.cn/7iR1kk",
				WxAppid:   "wxde8ac0a21135c07d",
				WxPath:    "/index/pages/h5/mtlm/mtlm?mt=3&lm=MTg2ODkzMzY2Mjg1MTkyNDA1MA%3D%3DNDY3%3D%3D%3D%3D&uid=85459&container=meituan_wxmini&lch=cps:x:0:65c5f4b9271221c79eae104d969a48a3:266983jutuike123456:408:85459",
			}, {
				Title:     "美团特惠酒店",
				Desc:      "美团特惠酒店在线预订，先领红包再下单。",
				ImageUrl:  "https://s3plus.sankuai.com/v1/mss_5017c592a8a946d2a54eb62a76ba299c/nebulafile/3967bcf477efa2d77819c819ebdbb3c9.png",
				TargetUrl: "https://kurl02.cn/7NV4GV",
				WxAppid:   "wxde8ac0a21135c07d",
				WxPath:    "/index/pages/h5/mtlm/mtlm?mt=3&lm=MTg2ODkzNDIxODMyNDM5ODEzMg%3D%3DNDY3%3D%3D%3D%3D&uid=85459&container=meituan_wxmini&lch=cps:x:0:65c5f4b9271221c79eae104d969a48a3:266983jutuike123456:400:85459",
			}, {
				Title:     "同程酒店优惠券",
				Desc:      "领百元入住红包，订房特惠5折起再减。",
				ImageUrl:  "https://img.jutuike.com/taokeout/banner/tcyl_hotal_banner.png?v=1",
				TargetUrl: "https://kzurl18.cn/7NVOib",
				WxAppid:   "wx336dcaf6a1ecf632",
				WxPath:    "page/home/webview/webview?src=https%3A%2F%2Fmp.elong.com%2Ftenthousandaura%2F%3Factivitycode%3D73086812-aaae-48ba-b14a-f087a6b61a92%26isSocket%3DHotel%26outToken%3D0d87b8e391ecbd01%26if%3D5012016%26of%3D5056928%26outUserid%3D10604_266983jutuike123456&isRefresh=refresh",
			}, {
				Title:     "飞猪酒店天天特惠",
				Desc:      "飞猪酒店天天特惠活动，最高立减150元，订房6折起。",
				ImageUrl:  "https://gw.alicdn.com/imgextra/i4/O1CN01vBpTLa1PVyJZJCxOS_!!6000000001847-2-tps-800-450.png",
				TargetUrl: "https://a.feizhu.com/3EkyQc",
				WxAppid:   "wx6a96c49f29850eb5",
				WxPath:    "pages/main/webview?url=https%3A%2F%2Foutfliggys.m.taobao.com%2Fxj%2Fpage%2Fpcraft%2Fpcraft%2Fcommon%2Fjiudianquanyi_copy21%3FtitleBarHidden%3D2%26ttid%3D12wechat000008258%26sht_track_info%3Dhw21arwxoc",
			},
		},
	}, {
		Label: "本地生活",
		Items: []common.WelfareItem{
			{
				Title:     "抖音团购",
				Desc:      "本地吃喝玩乐，低价爆品，上抖音一站购。",
				ImageUrl:  "https://img.jutuike.com/taokeout/poster/new/douyin_tuangou_banner.png",
				TargetUrl: "https://kurl11.cn/7NVaRE",
				WxAppid:   "wx89752980e795bfde",
				WxPath:    "/pages/plugin/index?type=tuangou&pub_id=266983&sid=123456&source=jutuike",
			},
		},
	},
}
