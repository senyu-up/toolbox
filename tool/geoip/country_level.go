package geoip

import "strings"

type Country struct {
	Id       string
	Code     int
	Language string
	Level    int
}

var confs map[string]*Country

func init() {
	confs = make(map[string]*Country, 302)

	confs["AE"] = &Country{"阿联酋", 101, "AE", 1}
	confs["AU"] = &Country{"澳大利亚", 107, "AU", 1}
	confs["CN"] = &Country{"中国", 127, "CN", 1}
	confs["JP"] = &Country{"日本", 166, "JP", 1}
	confs["KR"] = &Country{"韩国", 170, "KR", 1}
	confs["MC"] = &Country{"摩纳哥", 185, "MC", 1}
	confs["NZ"] = &Country{"新西兰", 205, "NZ", 1}
	confs["SG"] = &Country{"新加坡", 222, "SG", 1}
	confs["TW"] = &Country{"台湾省", 237, "TW", 1}
	confs["US"] = &Country{"美国", 242, "US", 1}
	confs["DE"] = &Country{"德国", 133, "DE", 1}
	confs["CA"] = &Country{"加拿大", 121, "CA", 1}
	confs["CH"] = &Country{"瑞士", 124, "CH", 1}
	confs["DK"] = &Country{"丹麦", 134, "DK", 1}
	confs["FR"] = &Country{"法国", 144, "FR", 1}
	confs["GB"] = &Country{"英国", 146, "GB", 1}
	confs["HK"] = &Country{"香港", 153, "HK", 1}
	confs["MO"] = &Country{"澳门", 191, "MO", 1}
	confs["QA"] = &Country{"卡塔尔", 215, "QA", 1}
	confs["SE"] = &Country{"瑞典", 221, "SE", 1}
	confs["KW"] = &Country{"科威特", 172, "KW", 2}
	confs["LU"] = &Country{"卢森堡", 181, "LU", 2}
	confs["NO"] = &Country{"挪威", 203, "NO", 2}
	confs["SA"] = &Country{"沙特阿拉伯", 218, "SA", 2}
	confs["BE"] = &Country{"比利时", 110, "BE", 2}
	confs["IS"] = &Country{"冰岛", 162, "IS", 2}
	confs["IT"] = &Country{"意大利", 163, "IT", 2}
	confs["LT"] = &Country{"立陶宛", 180, "LT", 2}
	confs["AT"] = &Country{"奥地利", 106, "AT", 2}
	confs["BN"] = &Country{"文莱", 116, "BN", 2}
	confs["BY"] = &Country{"白俄罗斯", 120, "BY", 2}
	confs["CL"] = &Country{"智利", 125, "CL", 2}
	confs["CS"] = &Country{"捷克", 130, "CS", 2}
	confs["EE"] = &Country{"爱沙尼亚", 138, "EE", 2}
	confs["ES"] = &Country{"西班牙", 140, "ES", 2}
	confs["FI"] = &Country{"芬兰", 142, "FI", 2}
	confs["GR"] = &Country{"希腊", 151, "GR", 2}
	confs["HU"] = &Country{"匈牙利", 155, "HU", 2}
	confs["IE"] = &Country{"爱尔兰", 157, "IE", 2}
	confs["IL"] = &Country{"以色列", 158, "IL", 2}
	confs["MX"] = &Country{"墨西哥", 195, "MX", 2}
	confs["NL"] = &Country{"荷兰", 202, "NL", 2}
	confs["PA"] = &Country{"巴拿马", 207, "PA", 2}
	confs["PL"] = &Country{"波兰", 212, "PL", 2}
	confs["PT"] = &Country{"葡萄牙", 213, "PT", 2}
	confs["RU"] = &Country{"俄罗斯", 217, "RU", 2}
	confs["TR"] = &Country{"土耳其", 236, "TR", 2}
	confs["UA"] = &Country{"乌克兰", 239, "UA", 2}
	confs["CY"] = &Country{"塞浦路斯", 132, "CY", 3}
	confs["KT"] = &Country{"科特迪瓦共和国", 171, "KT", 3}
	confs["LV"] = &Country{"拉脱维亚", 182, "LV", 3}
	confs["LY"] = &Country{"利比亚", 183, "LY", 3}
	confs["MA"] = &Country{"摩洛哥", 184, "MA", 3}
	confs["MD"] = &Country{"摩尔多瓦", 186, "MD", 3}
	confs["NP"] = &Country{"尼泊尔", 204, "NP", 3}
	confs["SI"] = &Country{"斯洛文尼亚", 223, "SI", 3}
	confs["SK"] = &Country{"斯洛伐克", 224, "SK", 3}
	confs["SM"] = &Country{"圣马力诺", 225, "SM", 3}
	confs["SN"] = &Country{"塞内加尔", 226, "SN", 3}
	confs["SO"] = &Country{"索马里", 227, "SO", 3}
	confs["UG"] = &Country{"乌干达", 240, "UG", 3}
	confs["UN"] = &Country{"联合国国旗", 241, "UN", 3}
	confs["BF"] = &Country{"布基纳法索", 111, "BF", 3}
	confs["BG"] = &Country{"保加利亚", 112, "BG", 3}
	confs["BH"] = &Country{"巴林", 113, "BH", 3}
	confs["BI"] = &Country{"布隆迪", 114, "BI", 3}
	confs["IQ"] = &Country{"伊拉克", 160, "IQ", 3}
	confs["IR"] = &Country{"伊朗", 161, "IR", 3}
	confs["LB"] = &Country{"黎巴嫩", 175, "LB", 3}
	confs["LC"] = &Country{"圣卢西亚", 176, "LC", 3}
	confs["LI"] = &Country{"列支敦士登", 177, "LI", 3}
	confs["LK"] = &Country{"斯里兰卡", 178, "LK", 3}
	confs["OM"] = &Country{"阿曼", 206, "OM", 3}
	confs["AF"] = &Country{"阿富汗", 102, "AF", 3}
	confs["AL"] = &Country{"阿尔巴尼亚", 103, "AL", 3}
	confs["AM"] = &Country{"亚美尼亚", 104, "AM", 3}
	confs["AR"] = &Country{"阿根廷", 105, "AR", 3}
	confs["AZ"] = &Country{"阿塞拜疆", 108, "AZ", 3}
	confs["BD"] = &Country{"孟加拉", 109, "BD", 3}
	confs["BJ"] = &Country{"贝宁", 115, "BJ", 3}
	confs["BO"] = &Country{"玻利维亚", 117, "BO", 3}
	confs["BR"] = &Country{"巴西", 118, "BR", 3}
	confs["BW"] = &Country{"博茨瓦纳", 119, "BW", 3}
	confs["CF"] = &Country{"中非", 122, "CF", 3}
	confs["CG"] = &Country{"刚果", 123, "CG", 3}
	confs["CM"] = &Country{"喀麦隆", 126, "CM", 3}
	confs["CO"] = &Country{"哥伦比亚", 128, "CO", 3}
	confs["CR"] = &Country{"哥斯达黎加", 129, "CR", 3}
	confs["CU"] = &Country{"古巴", 131, "CU", 3}
	confs["DO"] = &Country{"多米尼加共和国", 135, "DO", 3}
	confs["DZ"] = &Country{"阿尔及利亚", 136, "DZ", 3}
	confs["EC"] = &Country{"厄瓜多尔", 137, "EC", 3}
	confs["EG"] = &Country{"埃及", 139, "EG", 3}
	confs["ET"] = &Country{"埃塞俄比亚", 141, "ET", 3}
	confs["FJ"] = &Country{"斐济", 143, "FJ", 3}
	confs["GA"] = &Country{"加蓬", 145, "GA", 3}
	confs["GD"] = &Country{"格林纳达", 147, "GD", 3}
	confs["GE"] = &Country{"格鲁吉亚", 148, "GE", 3}
	confs["GH"] = &Country{"加纳", 149, "GH", 3}
	confs["GN"] = &Country{"几内亚", 150, "GN", 3}
	confs["GT"] = &Country{"危地马拉", 152, "GT", 3}
	confs["HN"] = &Country{"洪都拉斯", 154, "HN", 3}
	confs["ID"] = &Country{"印度尼西亚", 156, "ID", 3}
	confs["IN"] = &Country{"印度", 159, "IN", 3}
	confs["JM"] = &Country{"牙买加", 164, "JM", 3}
	confs["JO"] = &Country{"约旦", 165, "JO", 3}
	confs["KG"] = &Country{"吉尔吉斯坦", 167, "KG", 3}
	confs["KH"] = &Country{"柬埔寨", 168, "KH", 3}
	confs["KP"] = &Country{"北朝鲜", 169, "KP", 3}
	confs["KZ"] = &Country{"哈萨克", 173, "KZ", 3}
	confs["LA"] = &Country{"老挝", 174, "LA", 3}
	confs["LR"] = &Country{"利比里亚", 179, "LR", 3}
	confs["MG"] = &Country{"马达加斯加", 187, "MG", 3}
	confs["ML"] = &Country{"马里", 188, "ML", 3}
	confs["MM"] = &Country{"缅甸", 189, "MM", 3}
	confs["MN"] = &Country{"蒙古", 190, "MN", 3}
	confs["MT"] = &Country{"马耳他", 192, "MT", 3}
	confs["MU"] = &Country{"毛里求斯", 193, "MU", 3}
	confs["MW"] = &Country{"马拉维", 194, "MW", 3}
	confs["MY"] = &Country{"马来西亚", 196, "MY", 3}
	confs["MZ"] = &Country{"莫桑比克", 197, "MZ", 3}
	confs["NA"] = &Country{"纳米比亚", 198, "NA", 3}
	confs["NE"] = &Country{"尼日尔", 199, "NE", 3}
	confs["NG"] = &Country{"尼日利亚", 200, "NG", 3}
	confs["NI"] = &Country{"尼加拉瓜", 201, "NI", 3}
	confs["PE"] = &Country{"秘鲁", 208, "PE", 3}
	confs["PG"] = &Country{"巴布亚新几内亚", 209, "PG", 3}
	confs["PH"] = &Country{"菲律宾", 210, "PH", 3}
	confs["PK"] = &Country{"巴基斯坦", 211, "PK", 3}
	confs["PY"] = &Country{"巴拉圭", 214, "PY", 3}
	confs["RO"] = &Country{"罗马尼亚", 216, "RO", 3}
	confs["SC"] = &Country{"塞舌尔", 219, "SC", 3}
	confs["SD"] = &Country{"苏丹", 220, "SD", 3}
	confs["SY"] = &Country{"叙利亚", 228, "SY", 3}
	confs["SZ"] = &Country{"斯威士兰", 229, "SZ", 3}
	confs["TD"] = &Country{"乍得", 230, "TD", 3}
	confs["TG"] = &Country{"多哥", 231, "TG", 3}
	confs["TH"] = &Country{"泰国", 232, "TH", 3}
	confs["TJ"] = &Country{"塔吉克斯坦", 233, "TJ", 3}
	confs["TM"] = &Country{"土库曼", 234, "TM", 3}
	confs["TN"] = &Country{"突尼斯", 235, "TN", 3}
	confs["TZ"] = &Country{"坦桑尼亚", 238, "TZ", 3}
	confs["UY"] = &Country{"乌拉圭", 243, "UY", 3}
	confs["UZ"] = &Country{"乌兹别克", 244, "UZ", 3}
	confs["VC"] = &Country{"圣文森特岛", 245, "VC", 3}
	confs["VE"] = &Country{"委内瑞拉", 246, "VE", 3}
	confs["VN"] = &Country{"越南", 247, "VN", 3}
	confs["YE"] = &Country{"也门", 248, "YE", 3}
	confs["ZA"] = &Country{"南非", 249, "ZA", 3}
	confs["ZM"] = &Country{"赞比亚", 250, "ZM", 3}
	confs["ZW"] = &Country{"津巴布韦", 251, "ZW", 3}

	for code, _ := range confs {
		confs[confs[code].Id] = confs[code]
	}
}

func CountryConfigsByName(name string) *Country {
	conf, ok := confs[name]
	if ok {
		return conf
	}
	return nil
}

func CountryConfigsByCode(code string) *Country {
	conf, ok := confs[code]
	if ok {
		return conf
	}
	return nil
}

// 通过code获取国家
func CountryByCode(code string) string {
	conf, ok := confs[strings.ToUpper(code)]
	if !ok {
		return ""
	}
	if conf == nil {
		return ""
	}
	return conf.Id
}
