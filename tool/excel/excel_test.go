package excel

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

type Resume struct {
	Name       string    `json:"name" excel_name:"我的名字"`
	Like       string    `json:"like"`
	Sex        string    `json:"sex" excel_name:""`
	Age        int       `json:"age" excel_name:"我的年龄"`
	Status     int8      `json:"status" excel_name:"当前状态" enums:"0:冻结,1:开启,2:关闭"` //如设置enums 数据类型必须为 int8
	CreateTime int64     `json:"create_time" excel_name:"创建时间" excel_time:"int"`  //如设置 excel_time为int 数据类型必须为 int64
	UpdateTime time.Time `json:"update_time" excel_name:"更新时间" excel_time:"time"`
	Test       string    `json:"test" excel_name:"测试1"`
	Test2      string    `json:"test2" excel_name:"测试2"`
	Test3      string    `json:"test3" excel_name:"测试3"`
	Test4      string    `json:"test4" excel_name:"测试4"`
	Test5      string    `json:"test5" excel_name:"测试5"`
	Test6      string    `json:"test6" excel_name:"测试6"`
	Test7      string    `json:"test7" excel_name:"测试7"`
	Test8      string    `json:"test8" excel_name:"测试8"`
	Test9      string    `json:"test9" excel_name:"测试9"`
	Test15     string    `json:"test15" excel_name:"测试15"`
	Test16     string    `json:"test16" excel_name:"测试16"`
	Test17     string    `json:"test17" excel_name:"测试17"`
	Test18     string    `json:"test18" excel_name:"测试18"`
	Test19     string    `json:"test19" excel_name:"测试19"`
}

type HeatMapData struct {
	Index          string  `json:"index" excel_name:""`
	BuildingsPower float64 `json:"buildings_power" excel_name:"buildings_power"`
	Paid           float64 `json:"paid" excel_name:"paid"`
	Power          float64 `json:"power" excel_name:"power"`
	AllianceId     float64 `json:"alliance_id" excel_name:"alliance_id"`
	Energy         float64 `json:"energy" excel_name:"energy"`
	Language       float64 `json:"language" excel_name:"language"`
	TechsPower     float64 `json:"techs_power" excel_name:"techs_power"`
	PaidTimes      float64 `json:"paid_times" excel_name:"paid_times"`
	BaseLevel      float64 `json:"base_level" excel_name:"base_level"`
	Level          float64 `json:"level" excel_name:"level"`
	BanExpire      float64 `json:"ban_expire" excel_name:"ban_expire"`
	VipLevel       float64 `json:"vip_level" excel_name:"vip_level"`
	Avatar         float64 `json:"avatar" excel_name:"avatar"`
	ArmiesPower    float64 `json:"armies_power" excel_name:"armies_power"`
	IsInternal     float64 `json:"is_internal" excel_name:"is_internal"`
	LoginAt        float64 `json:"login_at" excel_name:"login_at"`
	LastLoginAt    float64 `json:"last_login_at" excel_name:"last_login_at"`
	CreatedAt      float64 `json:"created_at" excel_name:"created_at"`
	IsInitial      float64 `json:"is_initial" excel_name:"is_initial"`
	PaidAt         float64 `json:"paid_at" excel_name:"paid_at"`
	UpdateAt       float64 `json:"update_at" excel_name:"update_at"`
	LogoutAt       float64 `json:"logout_at" excel_name:"logout_at"`
}

func TestExcel(t *testing.T) {
	var stru []Resume
	info := Resume{
		Name:       "张三",
		Sex:        "男",
		Like:       "xxx",
		Age:        19,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now(),
		Status:     2,
		Test:       "1",
		Test2:      "1",
		Test3:      "1",
		Test4:      "1",
		Test5:      "1",
		Test6:      "1",
		Test7:      "1",
		Test8:      "1",
		Test9:      "1",
		Test15:     "1",
		Test16:     "1",
		Test17:     "1",
		Test18:     "1",
		Test19:     "1",
	}

	for i := 0; i < 10; i++ {
		stru = append(stru, info)
	}
	e := Excel{}
	err := e.SaveExcel("Book1.xlsx", &stru)
	fmt.Println(err)
	return
}
func TestCsvReadExcel(t *testing.T) {
	file, err := os.Open("result.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(record) // record has the type []string
	}

	return
}

type Category struct {
	Q1 string `json:"q_1" excel_name:"一级"`
	Q2 string `json:"q_2" excel_name:"二级"`
	Q3 string `json:"q_3" excel_name:"三级"`
	Q4 string `json:"q_4" excel_name:"四级"`
}

func TestExcel_ReadExcel(t *testing.T) {
	e := Excel{
		Sheet: "问题类型",
	}
	var stru []*Category
	err := e.ReadExcel("问题类型.xlsx", &stru)
	fmt.Println(err)
	for i, v := range stru {
		if v.Q1 == "" {
			stru[i].Q1 = stru[i-1].Q1
		}
		if v.Q2 == "" {
			stru[i].Q2 = stru[i-1].Q2
		}
		if v.Q3 == "" {
			stru[i].Q3 = stru[i-1].Q3
		}
		if v.Q4 == "" {
			stru[i].Q4 = stru[i-1].Q4
		}
	}
	jsonStr, _ := json.Marshal(stru)
	fmt.Println(string(jsonStr))
	return
}

type chartDataRow struct {
	X     string  `json:"x"`
	Y     string  `json:"y"`
	Value float64 `json:"value"`
}
type chartData struct {
	Columns []string       `json:"columns"`
	Rows    []chartDataRow `json:"rows"`
}

func TestExcel_ReadCsv(t *testing.T) {
	e := Excel{}
	var stru []*HeatMapData
	err := e.ReadCsv("result.csv", &stru)
	fmt.Println(err)
	jsonStr, _ := json.Marshal(stru)
	fmt.Println(string(jsonStr))
	ChartData := chartData{}
	ChartData.Columns = []string{"x", "y", "value"}
	for _, v := range stru {
		info := &chartDataRow{}
		vStr, _ := json.Marshal(v)
		vMap := make(map[string]interface{})
		_ = json.Unmarshal(vStr, &vMap)
		info.X = v.Index
		for k, e := range vMap {
			if k == "index" {
				continue
			}
			info.Y = k
			info.Value = e.(float64)
			ChartData.Rows = append(ChartData.Rows, *info)
		}
	}
	jsonStr, _ = json.Marshal(ChartData)
	fmt.Println(string(jsonStr))
	return
}

func TestRunPythonScript(t *testing.T) {
	out, err := exec.Command("python3", "heatmap.py").Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	result := string(out)
	fmt.Println("result:", result)
	if strings.Index(result, "success") == -1 {
		err = errors.New(fmt.Sprintf("error：%s", result))
	}
	fmt.Println("err:", err)
	return
}
func TestAddStr(t *testing.T) {

	fmt.Println(addStr("BZ", 3))
}

func Test360ReadExcel(t *testing.T) {
	f, err := excelize.OpenFile("result.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Get value from cell by given worksheet name and axis.
	cell, err := f.GetCellValue("result", "B2")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cell)
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("result")
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}
func TestWriteCsv(t *testing.T) {
	var stru []Resume
	info := Resume{
		Name:       "张三",
		Sex:        "男",
		Like:       "xxx",
		Age:        19,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now(),
		Status:     2,
		Test:       "1",
		Test2:      "1",
		Test3:      "1",
		Test4:      "1",
		Test5:      "1",
		Test6:      "1",
		Test7:      "1",
		Test8:      "1",
		Test9:      "1",
		Test15:     "1",
		Test16:     "1",
		Test17:     "1",
		Test18:     "1",
		Test19:     "1",
	}

	for i := 0; i < 10; i++ {
		stru = append(stru, info)
	}
	e := Excel{}
	err := e.SaveCsv("test_02.csv", &stru)
	fmt.Println(err)
	return
}

func TestMapSaveCsv(t *testing.T) {
	data := make([]map[string]interface{}, 0)
	str := `[
			{
                "ac_id": 0,
                "acc_id": 1,
                "action": 2,
                "action_id": 3,
                "balance": 4
            },
			{
                "ac_id": "活动id",
                "acc_id": "uid",
                "action": "操作信息",
                "action_id": "操作id",
                "balance": "余额"
            },
            {
                "ac_id": "2010501",
                "acc_id": "1000011671",
                "action": "个人全局战力单服积分",
                "action_id": "c0egn9bpg1fpsahfo5k0",
                "balance": "73167",
                "created_at": "1612516261",
                "day": "5",
                "extend_1": "",
                "extend_2": "",
                "extend_3": "",
                "extend_4": "",
                "extend_5": "",
                "month": "2",
                "nums": "11",
                "pk_id": "1357617716038467584",
                "role_id": "S2#4344",
                "role_name": "Guest#6017DE04",
                "save_time": "1612516266",
                "server_id": "S2",
                "year": "2021"
            },
            {
                "ac_id": "2010501",
                "acc_id": "1000011664",
                "action": "个人全局战力单服积分",
                "action_id": "c0egn9bpg1fpsahfo5m0",
                "balance": "73167",
                "created_at": "1612516261",
                "day": "5",
                "extend_1": "",
                "extend_2": "",
                "extend_3": "",
                "extend_4": "",
                "extend_5": "",
                "month": "2",
                "nums": "11",
                "pk_id": "1357617717296758784",
                "role_id": "S2#4341",
                "role_name": "Guest#6017DDFE",
                "save_time": "1612516266",
                "server_id": "S2",
                "year": "2021"
            },
            {
                "ac_id": "2010501",
                "acc_id": "1000011576",
                "action": "个人全局战力单服积分",
                "action_id": "c0egn93pg1fpsahfo5f0",
                "balance": "114893",
                "created_at": "1612516260",
                "day": "5",
                "extend_1": "",
                "extend_2": "",
                "extend_3": "",
                "extend_4": "",
                "extend_5": "",
                "month": "2",
                "nums": "525",
                "pk_id": "1357617710262910976",
                "role_id": "S2#4314",
                "role_name": "Julius",
                "save_time": "1612516266",
                "server_id": "S2",
                "year": "2021"
            },
            {
                "ac_id": "2010501",
                "acc_id": "1000011576",
                "action": "个人全局战力单服积分",
                "action_id": "c0egn93pg1fpsahfo5i0",
                "balance": "115943",
                "created_at": "1612516260",
                "day": "5",
                "extend_1": "",
                "extend_2": "",
                "extend_3": "",
                "extend_4": "",
                "extend_5": "",
                "month": "2",
                "nums": "525",
                "pk_id": "1357617711575728128",
                "role_id": "S2#4314",
                "role_name": "Julius",
                "save_time": "1612516266",
                "server_id": "S2",
                "year": "2021"
            },
            {
                "ac_id": "2010501",
                "acc_id": "1000011576",
                "action": "个人全局战力单服积分",
                "action_id": "c0egn93pg1fpsahfo5gg",
                "balance": "115418",
                "created_at": "1612516260",
                "day": "5",
                "extend_1": "",
                "extend_2": "",
                "extend_3": "",
                "extend_4": "",
                "extend_5": "",
                "month": "2",
                "nums": "525",
                "pk_id": "1357617710862696448",
                "role_id": "S2#4314",
                "role_name": "Julius",
                "save_time": "1612516266",
                "server_id": "S2",
                "year": "2021"
            },
            {
                "ac_id": "2010501",
                "acc_id": "1000011576",
                "action": "个人全局战力单服积分",
                "action_id": "c0egn8bpg1fpsahfo5bg",
                "balance": "114368",
                "created_at": "1612516257",
                "day": "5",
                "extend_1": "",
                "extend_2": "",
                "extend_3": "",
                "extend_4": "",
                "extend_5": "",
                "month": "2",
                "nums": "5775",
                "pk_id": "1357617700284661760",
                "role_id": "S2#4314",
                "role_name": "Julius",
                "save_time": "1612516266",
                "server_id": "S2",
                "year": "2021"
            },
            {
                "ac_id": "2010501",
                "acc_id": "1000011576",
                "action": "个人全局战力单服积分",
                "action_id": "c0egn7jpg1fpsahfo55g",
                "balance": "107543",
                "created_at": "1612516254",
                "day": "5",
                "extend_1": "",
                "extend_2": "",
                "extend_3": "",
                "extend_4": "",
                "extend_5": "",
                "month": "2",
                "nums": "525",
                "pk_id": "1357617687244570624",
                "role_id": "S2#4314",
                "role_name": "Julius",
                "save_time": "1612516266",
                "server_id": "S2",
                "year": "2021"
            },
            {
                "ac_id": "2010501",
                "acc_id": "1000011576",
                "action": "个人全局战力单服积分",
                "action_id": "c0egn7jpg1fpsahfo53g",
                "balance": "107018",
                "created_at": "1612516254",
                "day": "5",
                "extend_1": "",
                "extend_2": "",
                "extend_3": "",
                "extend_4": "",
                "extend_5": "",
                "month": "2",
                "nums": "525",
                "pk_id": "1357617686275686400",
                "role_id": "S2#4314",
                "role_name": "Julius",
                "save_time": "1612516266",
                "server_id": "S2",
                "year": "2021"
            },
            {
                "ac_id": "2010501",
                "acc_id": "1000011576",
                "action": "个人全局战力单服积分",
                "action_id": "c0egn7jpg1fpsahfo570",
                "balance": "108068",
                "created_at": "1612516254",
                "day": "5",
                "extend_1": "",
                "extend_2": "",
                "extend_3": "",
                "extend_4": "",
                "extend_5": "",
                "month": "2",
                "nums": "525",
                "pk_id": "1357617687756275712",
                "role_id": "S2#4314",
                "role_name": "Julius",
                "save_time": "1612516266",
                "server_id": "S2",
                "year": "2021"
            },
            {
                "ac_id": "2010501",
                "acc_id": "1000011576",
                "action": "个人全局战力单服积分",
                "action_id": "c0egn7jpg1fpsahfo58g",
                "balance": "108593",
                "created_at": "1612516254",
                "day": "5",
                "extend_1": "",
                "extend_2": "",
                "extend_3": "",
                "extend_4": "",
                "extend_5": "",
                "month": "2",
                "nums": "525",
                "pk_id": "1357617688297340928",
                "role_id": "S2#4314",
                "role_name": "Julius",
                "save_time": "1612516266",
                "server_id": "S2",
                "year": "2021"
            },
            {
                "ac_id": "2010501",
                "acc_id": "1000008458",
                "action": "个人全局战力单服积分",
                "action_id": "c0egn63pg1fpsahfo4sg",
                "balance": "360645",
                "created_at": "1612516248",
                "day": "5",
                "extend_1": "",
                "extend_2": "",
                "extend_3": "",
                "extend_4": "",
                "extend_5": "",
                "month": "2",
                "nums": "570",
                "pk_id": "1357617662372347904",
                "role_id": "S2#1480",
                "role_name": "Guest#6017C7BA",
                "save_time": "1612516253",
                "server_id": "S2",
                "year": "2021"
            },
            {
                "ac_id": "2010501",
                "acc_id": "1000011576",
                "action": "个人全局战力单服积分",
                "action_id": "c0egn5bpg1fpsahfo4m0",
                "balance": "106493",
                "created_at": "1612516245",
                "day": "5",
                "extend_1": "",
                "extend_2": "",
                "extend_3": "",
                "extend_4": "",
                "extend_5": "",
                "month": "2",
                "nums": "11",
                "pk_id": "1357617648279486464",
                "role_id": "S2#4314",
                "role_name": "Julius",
                "save_time": "1612516253",
                "server_id": "S2",
                "year": "2021"
            },
            {
                "ac_id": "2010501",
                "acc_id": "1000011576",
                "action": "个人全局战力单服积分",
                "action_id": "c0egn3bpg1fpsahfo4a0",
                "balance": "106482",
                "created_at": "1612516237",
                "day": "5",
                "extend_1": "",
                "extend_2": "",
                "extend_3": "",
                "extend_4": "",
                "extend_5": "",
                "month": "2",
                "nums": "11",
                "pk_id": "1357617614670528512",
                "role_id": "S2#4314",
                "role_name": "Julius",
                "save_time": "1612516242",
                "server_id": "S2",
                "year": "2021"
            },
            {
                "ac_id": "2010501",
                "acc_id": "1000011671",
                "action": "个人全局战力单服积分",
                "action_id": "c0egn3bpg1fpsahfo4dg",
                "balance": "73156",
                "created_at": "1612516237",
                "day": "5",
                "extend_1": "",
                "extend_2": "",
                "extend_3": "",
                "extend_4": "",
                "extend_5": "",
                "month": "2",
                "nums": "73156",
                "pk_id": "1357617615819767808",
                "role_id": "S2#4344",
                "role_name": "Guest#6017DE04",
                "save_time": "1612516242",
                "server_id": "S2",
                "year": "2021"
            },
            {
                "ac_id": "2010501",
                "acc_id": "1000008458",
                "action": "个人全局战力单服积分",
                "action_id": "c0egn2rpg1fpsahfo47g",
                "balance": "360075",
                "created_at": "1612516235",
                "day": "5",
                "extend_1": "",
                "extend_2": "",
                "extend_3": "",
                "extend_4": "",
                "extend_5": "",
                "month": "2",
                "nums": "570",
                "pk_id": "1357617608769142784",
                "role_id": "S2#1480",
                "role_name": "Guest#6017C7BA",
                "save_time": "1612516242",
                "server_id": "S2",
                "year": "2021"
            },
            {
                "ac_id": "2010501",
                "acc_id": "1000011580",
                "action": "个人全局战力单服积分",
                "action_id": "c0egn2jpg1fpsahfo450",
                "balance": "104711",
                "created_at": "1612516234",
                "day": "5",
                "extend_1": "",
                "extend_2": "",
                "extend_3": "",
                "extend_4": "",
                "extend_5": "",
                "month": "2",
                "nums": "55",
                "pk_id": "1357617601986953216",
                "role_id": "S2#4317",
                "role_name": "Guest#6017DDD0",
                "save_time": "1612516242",
                "server_id": "S2",
                "year": "2021"
            },
            {
                "ac_id": "2010501",
                "acc_id": "1000011664",
                "action": "个人全局战力单服积分",
                "action_id": "c0egn2bpg1fpsahfo42g",
                "balance": "73156",
                "created_at": "1612516233",
                "day": "5",
                "extend_1": "",
                "extend_2": "",
                "extend_3": "",
                "extend_4": "",
                "extend_5": "",
                "month": "2",
                "nums": "73156",
                "pk_id": "1357617599797526528",
                "role_id": "S2#4341",
                "role_name": "Guest#6017DDFE",
                "save_time": "1612516242",
                "server_id": "S2",
                "year": "2021"
            },
            {
                "ac_id": "2010501",
                "acc_id": "1000011171",
                "action": "个人全局战力单服积分",
                "action_id": "c0egmvbpg1fpsahfo3h0",
                "balance": "147135",
                "created_at": "1612516221",
                "day": "5",
                "extend_1": "",
                "extend_2": "",
                "extend_3": "",
                "extend_4": "",
                "extend_5": "",
                "month": "2",
                "nums": "43",
                "pk_id": "1357617547435835392",
                "role_id": "S2#4100",
                "role_name": "white devil",
                "save_time": "1612516226",
                "server_id": "S2",
                "year": "2021"
            },
            {
                "ac_id": "2010501",
                "acc_id": "1000008458",
                "action": "个人全局战力单服积分",
                "action_id": "c0egmv3pg1fpsahfo3eg",
                "balance": "359505",
                "created_at": "1612516220",
                "day": "5",
                "extend_1": "",
                "extend_2": "",
                "extend_3": "",
                "extend_4": "",
                "extend_5": "",
                "month": "2",
                "nums": "570",
                "pk_id": "1357617542281035776",
                "role_id": "S2#1480",
                "role_name": "Guest#6017C7BA",
                "save_time": "1612516226",
                "server_id": "S2",
                "year": "2021"
            },
            {
                "ac_id": "1900601",
                "acc_id": "1000011664",
                "action": "七日活动领取奖励",
                "action_id": "c0egmrbpg1fpsahfo2p0",
                "balance": "5",
                "created_at": "1612516205",
                "day": "5",
                "extend_1": "",
                "extend_2": "",
                "extend_3": "",
                "extend_4": "",
                "extend_5": "",
                "month": "2",
                "nums": "5",
                "pk_id": "1357617480066924544",
                "role_id": "S2#4341",
                "role_name": "Guest#6017DDFE",
                "save_time": "1612516210",
                "server_id": "S2",
                "year": "2021"
            }
        ]`
	err := json.Unmarshal([]byte(str), &data)
	fmt.Println(err)
	e := Excel{}
	err = e.SaveExcel("test_02.xlsx", data)
	fmt.Println(err)
}
