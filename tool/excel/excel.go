package excel

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	filetool "github.com/senyu-up/toolbox/tool/file"
	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
)

type Excel struct {
	f     *excelize.File
	Sheet string
}

func (e *Excel) ReadExcel(file string, data interface{}) (err error) {
	if reflect.ValueOf(data).Kind() != reflect.Ptr || reflect.TypeOf(data).Elem().Kind() != reflect.Slice {
		err = errors.New("参数错误")
		return err
	}
	e.f, err = excelize.OpenFile(file)
	if err != nil {
		return err
	}
	v := reflect.ValueOf(data).Elem()
	t := reflect.TypeOf(data).Elem().Elem()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	//如果传入的excel页为空，就获取一页sheet
	if e.Sheet == "" {
		e.Sheet = e.f.GetSheetName(0)
	}
	// Get all the rows in the Sheet1.
	rows, err := e.f.GetRows(e.Sheet)
	var map_data = make(map[int]int)
	var map_arr = make(map[int]map[string]int)
	var v_index = 0
	for j, row := range rows {
		if j == 0 {
			for y, colCell := range row {
				if j == 0 {
					null_num := 0
					for i := 0; i < t.NumField(); i++ {
						excel_name, ok := t.Field(i).Tag.Lookup("excel_name")
						if !ok {
							null_num++
							continue
						}
						if excel_name == colCell {
							map_data[i] = y
						}
						if enums_str := t.Field(i).Tag.Get("enums"); enums_str != "" {
							var map_v = make(map[string]int)
							enums_arr := strings.Split(enums_str, ",")
							for _, v := range enums_arr {
								enums_v := strings.Split(v, ":")
								key, _ := strconv.Atoi(enums_v[0])
								map_v[enums_v[1]] = key
							}
							map_arr[i] = map_v
						}
					}
				}

			}
			continue
		}
		var subv reflect.Value
		subv = reflect.New(t)
		if v.Type().Elem().Kind() != reflect.Ptr {
			subv = subv.Elem()
		}
		v2 := reflect.Append(v, subv)
		v.Set(v2)
		v_index_value := v.Index(v_index)
		if v_index_value.Type().Kind() == reflect.Ptr {
			v_index_value = v_index_value.Elem()
		}
		for i := 0; i < v_index_value.NumField(); i++ {
			if _, ok := map_data[i]; !ok {
				continue
			}
			excel_value := ""
			if len(row) > map_data[i] {
				excel_value = row[map_data[i]]
			}
			switch v_index_value.Field(i).Type().Kind() {
			case reflect.String:
				v_index_value.Field(i).Set(reflect.ValueOf(excel_value))
				break
			case reflect.Int:
				data_v, _ := strconv.Atoi(excel_value)
				v_index_value.Field(i).Set(reflect.ValueOf(data_v))
				break
			case reflect.Int8:
				data_v, _ := strconv.Atoi(excel_value)
				if enums_str := t.Field(i).Tag.Get("enums"); enums_str != "" {
					map_v := map_arr[i]
					if _, ok := map_v[excel_value]; ok {
						data_v = map_v[excel_value]
					}
				}
				v_index_value.Field(i).Set(reflect.ValueOf(int8(data_v)))
				break
			case reflect.Int16:
				data_v, _ := strconv.Atoi(excel_value)
				v_index_value.Field(i).Set(reflect.ValueOf(int16(data_v)))
				break
			case reflect.Int32:
				data_v, _ := strconv.Atoi(excel_value)
				v_index_value.Field(i).Set(reflect.ValueOf(int32(data_v)))
				break
			case reflect.Int64:
				data_v, _ := strconv.Atoi(excel_value)
				if t.Field(i).Tag.Get("excel_time") == "int" || t.Field(i).Tag.Get("excel_time") == "int64" {
					local, _ := time.LoadLocation("Local")
					t, _ := time.ParseInLocation("2006-01-02 15:04:05", excel_value, local)
					data_v = int(t.Unix())
				}
				v_index_value.Field(i).Set(reflect.ValueOf(int64(data_v)))
				break
			case reflect.Uint:
				data_v, _ := strconv.Atoi(excel_value)
				v_index_value.Field(i).Set(reflect.ValueOf(uint(data_v)))
				break
			case reflect.Uint8:
				data_v, _ := strconv.Atoi(excel_value)
				v_index_value.Field(i).Set(reflect.ValueOf(uint8(data_v)))
				break
			case reflect.Uint16:
				data_v, _ := strconv.Atoi(excel_value)
				v_index_value.Field(i).Set(reflect.ValueOf(uint16(data_v)))
				break
			case reflect.Uint32:
				data_v, _ := strconv.Atoi(excel_value)
				v_index_value.Field(i).Set(reflect.ValueOf(uint32(data_v)))
				break
			case reflect.Uint64:
				data_v, _ := strconv.Atoi(excel_value)
				v_index_value.Field(i).Set(reflect.ValueOf(uint64(data_v)))
				break
			case reflect.Float32:
				data_v, _ := strconv.ParseFloat(excel_value, 64)
				v_index_value.Field(i).Set(reflect.ValueOf(float32(data_v)))
				break
			case reflect.Float64:
				data_v, _ := strconv.ParseFloat(excel_value, 64)
				v_index_value.Field(i).Set(reflect.ValueOf(data_v))
				break
			case reflect.Bool:
				data_v, _ := strconv.ParseBool(excel_value)
				v_index_value.Field(i).Set(reflect.ValueOf(data_v))
				break
			case reflect.TypeOf(time.Time{}).Kind():
				local, _ := time.LoadLocation("Local")
				t, _ := time.ParseInLocation("2006-01-02 15:04:05", excel_value, local)
				v_index_value.Field(i).Set(reflect.ValueOf(t))
				break
			default:
				//v_index_value.Field(i).Set(reflect.ValueOf(0))
				//break
			}
		}
		v_index++

	}
	return
}
func (e *Excel) ReadCsv(file string, data interface{}) (err error) {
	if reflect.ValueOf(data).Kind() != reflect.Ptr || reflect.TypeOf(data).Elem().Kind() != reflect.Slice {
		err = errors.New("参数错误")
		return err
	}
	//e.f ,err = excelize.OpenFile(file)
	//if err != nil {
	//	return err
	//}
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	reader := csv.NewReader(f)
	v := reflect.ValueOf(data).Elem()
	t := reflect.TypeOf(data).Elem().Elem()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// Get all the rows in the Sheet1.
	var map_data = make(map[int]int)
	var map_arr = make(map[int]map[string]int)
	var v_index = 0
	var j = 0
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if j == 0 {
			for y, colCell := range row {
				null_num := 0
				for i := 0; i < t.NumField(); i++ {
					excel_name, ok := t.Field(i).Tag.Lookup("excel_name")
					if !ok {
						null_num++
						continue
					}
					if excel_name == colCell {
						map_data[i] = y
					}
					if enums_str := t.Field(i).Tag.Get("enums"); enums_str != "" {
						var map_v = make(map[string]int)
						enums_arr := strings.Split(enums_str, ",")
						for _, v := range enums_arr {
							enums_v := strings.Split(v, ":")
							key, _ := strconv.Atoi(enums_v[0])
							map_v[enums_v[1]] = key
						}
						map_arr[i] = map_v
					}
				}

			}
			j++
			continue
		}
		var subv reflect.Value
		subv = reflect.New(t)
		if v.Type().Elem().Kind() != reflect.Ptr {
			subv = subv.Elem()
		}
		v2 := reflect.Append(v, subv)
		v.Set(v2)
		v_index_value := v.Index(v_index)
		if v_index_value.Type().Kind() == reflect.Ptr {
			v_index_value = v_index_value.Elem()
		}
		for i := 0; i < v_index_value.NumField(); i++ {
			if _, ok := map_data[i]; !ok {
				continue
			}
			excel_value := ""
			if len(row) > map_data[i] {
				excel_value = row[map_data[i]]
			}
			switch v_index_value.Field(i).Type().Kind() {
			case reflect.String:
				v_index_value.Field(i).Set(reflect.ValueOf(excel_value))
				break
			case reflect.Int:
				data_v, _ := strconv.Atoi(excel_value)
				v_index_value.Field(i).Set(reflect.ValueOf(data_v))
				break
			case reflect.Int8:
				data_v, _ := strconv.Atoi(excel_value)
				if enums_str := t.Field(i).Tag.Get("enums"); enums_str != "" {
					map_v := map_arr[i]
					if _, ok := map_v[excel_value]; ok {
						data_v = map_v[excel_value]
					}
				}
				v_index_value.Field(i).Set(reflect.ValueOf(int8(data_v)))
				break
			case reflect.Int16:
				data_v, _ := strconv.Atoi(excel_value)
				v_index_value.Field(i).Set(reflect.ValueOf(int16(data_v)))
				break
			case reflect.Int32:
				data_v, _ := strconv.Atoi(excel_value)
				v_index_value.Field(i).Set(reflect.ValueOf(int32(data_v)))
				break
			case reflect.Int64:
				data_v, _ := strconv.Atoi(excel_value)
				if t.Field(i).Tag.Get("excel_time") == "int" || t.Field(i).Tag.Get("excel_time") == "int64" {
					local, _ := time.LoadLocation("Local")
					t, _ := time.ParseInLocation("2006-01-02 15:04:05", excel_value, local)
					data_v = int(t.Unix())
				}
				v_index_value.Field(i).Set(reflect.ValueOf(int64(data_v)))
				break
			case reflect.Uint:
				data_v, _ := strconv.Atoi(excel_value)
				v_index_value.Field(i).Set(reflect.ValueOf(uint(data_v)))
				break
			case reflect.Uint8:
				data_v, _ := strconv.Atoi(excel_value)
				v_index_value.Field(i).Set(reflect.ValueOf(uint8(data_v)))
				break
			case reflect.Uint16:
				data_v, _ := strconv.Atoi(excel_value)
				v_index_value.Field(i).Set(reflect.ValueOf(uint16(data_v)))
				break
			case reflect.Uint32:
				data_v, _ := strconv.Atoi(excel_value)
				v_index_value.Field(i).Set(reflect.ValueOf(uint32(data_v)))
				break
			case reflect.Uint64:
				data_v, _ := strconv.Atoi(excel_value)
				v_index_value.Field(i).Set(reflect.ValueOf(uint64(data_v)))
				break
			case reflect.Float32:
				data_v, _ := strconv.ParseFloat(excel_value, 64)
				v_index_value.Field(i).Set(reflect.ValueOf(float32(data_v)))
				break
			case reflect.Float64:
				data_v, _ := strconv.ParseFloat(excel_value, 64)
				v_index_value.Field(i).Set(reflect.ValueOf(data_v))
				break
			case reflect.Bool:
				data_v, _ := strconv.ParseBool(excel_value)
				v_index_value.Field(i).Set(reflect.ValueOf(data_v))
				break
			case reflect.TypeOf(time.Time{}).Kind():
				local, _ := time.LoadLocation("Local")
				t, _ := time.ParseInLocation("2006-01-02 15:04:05", excel_value, local)
				v_index_value.Field(i).Set(reflect.ValueOf(t))
				break
			default:
				//v_index_value.Field(i).Set(reflect.ValueOf(0))
				//break
			}
		}
		v_index++
	}
	return
}
func (e *Excel) GetRowsNum() (num int) {
	if e.f == nil {
		return 0
	}
	rows, err := e.f.Rows(e.Sheet)
	if err != nil {
		return 0
	}
	for rows.Next() {
		num++
	}
	return num
}
func (e *Excel) SaveExcel(file string, data interface{}) (err error) {
	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		if reflect.TypeOf(data).Elem().Kind() == reflect.Map {
			if dataMap, ok := data.([]map[string]interface{}); ok {
				return e.MapSaveExcel(file, dataMap)
			}
			err = errors.New("参数错误")
			return err
		}
	} else {
		if reflect.TypeOf(data).Elem().Elem().Kind() == reflect.Map {
			if dataMap, ok := data.(*[]map[string]interface{}); ok {
				return e.MapSaveExcel(file, *dataMap)
			}
			err = errors.New("参数错误")
			return err
		}
	}
	if e.Sheet == "" {
		e.Sheet = "Sheet1"
	}
	if reflect.ValueOf(data).Kind() != reflect.Ptr || reflect.TypeOf(data).Elem().Kind() != reflect.Slice {
		err = errors.New("参数错误")
		return err
	}
	e.f = excelize.NewFile()
	e.f, err = excelize.OpenFile(file)
	var index int
	var beginRow = 1
	if err != nil {
		e.f = excelize.NewFile()
		index, _ = e.f.NewSheet(e.Sheet)
	} else {
		index, _ = e.f.GetSheetIndex(e.Sheet)
		beginRow = e.GetRowsNum()
	}
	t := reflect.TypeOf(data).Elem().Elem()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	null_num := 0
	var map_arr = make(map[int]map[int]string)
	for i := 0; i < t.NumField(); i++ {
		excel_name, ok := t.Field(i).Tag.Lookup("excel_name")
		if !ok {
			null_num++
			continue
		}
		enums_str := t.Field(i).Tag.Get("enums")
		if enums_str != "" {
			map_data := make(map[int]string)
			enums_arr := strings.Split(enums_str, ",")
			for _, v := range enums_arr {
				enums_v := strings.Split(v, ":")
				key, _ := strconv.Atoi(enums_v[0])
				map_data[key] = enums_v[1]
			}
			map_arr[i] = map_data
		}
		if beginRow > 1 {
			continue
		}
		axis := fmt.Sprintf("%s%d", addStr("A", int32(i-null_num)), beginRow)
		e.f.SetCellValue(e.Sheet, axis, excel_name)
	}
	num := reflect.ValueOf(data).Elem().Len()
	if num > 0 {
		for i := 0; i < num; i++ {
			value := reflect.ValueOf(data).Elem().Index(i)
			if reflect.ValueOf(data).Elem().Index(i).Type().Kind() == reflect.Ptr {
				value = value.Elem()
			}
			t := value.Type()
			null_num := 0
			for x := 0; x < t.NumField(); x++ {
				_, ok := t.Field(x).Tag.Lookup("excel_name")
				if !ok {
					null_num++
					continue
				}
				cell_value := value.Field(x)
				enums_str := t.Field(x).Tag.Get("enums")
				if enums_str != "" {
					map_data := map_arr[x]
					if _, ok := map_data[int(value.Field(x).Int())]; ok {
						cell_value = reflect.ValueOf(map_data[int(value.Field(x).Int())])
					}
				}
				excel_time := t.Field(x).Tag.Get("excel_time")
				switch excel_time {
				case "time":
					val_str := fmt.Sprintf("%v", value.Field(x))
					cell_value = reflect.ValueOf(val_str[:19])
					break
				case "int", "int64":
					t := time.Unix(value.Field(x).Int(), 0).Format("2006-01-02 15:04:05")
					cell_value = reflect.ValueOf(t)
					break
				}
				axis := fmt.Sprintf("%s%d", addStr("A", int32(x-null_num)), i+1+beginRow)
				e.f.SetCellValue(e.Sheet, axis, cell_value)
			}
		}
	}
	// Set value of a cell.
	// Set active sheet of the workbook.
	e.f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	err = e.f.SaveAs(file)
	return
}
func (e *Excel) MapSaveExcel(file string, data []map[string]interface{}) (err error) {
	if len(data) < 2 {
		return errors.New("数据有问题")
	}
	if e.Sheet == "" {
		e.Sheet = "Sheet1"
	}
	e.f = excelize.NewFile()
	e.f, err = excelize.OpenFile(file)
	var index int
	var beginRow = 1
	if err != nil {
		e.f = excelize.NewFile()
		index, _ = e.f.NewSheet(e.Sheet)
	} else {
		index, _ = e.f.GetSheetIndex(e.Sheet)
		beginRow = e.GetRowsNum()
	}
	t := reflect.TypeOf(data).Elem().Elem()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	dataSort := data[0]
	csvHeard := make([]string, len(dataSort))
	for k, v := range data[1] {
		if sort, ok := dataSort[k]; ok {
			if cast.ToInt(sort) > len(csvHeard)-1 {
				return errors.New("表头错误")
			}
			csvHeard[cast.ToInt(sort)] = fmt.Sprintf("%v", v)
		}
	}
	for i, heard := range csvHeard {
		if beginRow > 1 {
			continue
		}
		axis := fmt.Sprintf("%s%d", addStr("A", int32(i)), beginRow)
		e.f.SetCellValue(e.Sheet, axis, heard)
	}
	tempData := data[2:]
	for tempI, temp := range tempData {
		cellData := make([]string, len(csvHeard))
		for k, v := range temp {
			if sort, ok := dataSort[k]; ok {
				cellData[cast.ToInt(sort)] = fmt.Sprintf("%v", v)
			}
		}
		for i, value := range cellData {
			axis := fmt.Sprintf("%s%d", addStr("A", int32(i)), tempI+1+beginRow)
			e.f.SetCellValue(e.Sheet, axis, value)
		}
	}
	// Set value of a cell.
	// Set active sheet of the workbook.
	e.f.SetActiveSheet(index)
	// Save spreadsheet by the given path.
	err = e.f.SaveAs(file)
	return
}

func (e *Excel) SaveCsv(file string, data interface{}) (err error) {
	if reflect.ValueOf(data).Kind() != reflect.Ptr {
		if reflect.TypeOf(data).Elem().Kind() == reflect.Map {
			if dataMap, ok := data.([]map[string]interface{}); ok {
				return e.MapSaveCsv(file, dataMap)
			}
			err = errors.New("参数错误")
			return err
		}
	} else {
		if reflect.TypeOf(data).Elem().Elem().Kind() == reflect.Map {
			if dataMap, ok := data.(*[]map[string]interface{}); ok {
				return e.MapSaveCsv(file, *dataMap)
			}
			err = errors.New("参数错误")
			return err
		}
	}
	if reflect.ValueOf(data).Kind() != reflect.Ptr || reflect.TypeOf(data).Elem().Kind() != reflect.Slice {
		err = errors.New("参数错误")
		return err
	}
	alread := 1
	f := &os.File{}
	if filetool.FileIsExist(file) {
		f, err = os.OpenFile(file, os.O_WRONLY|os.O_APPEND, 0666)
	} else {
		f, err = os.Create(file)
		if err != nil {
			return err
		}
		alread = 0
	}
	defer f.Close()
	if alread == 0 {
		// 写入UTF-8 BOM，防止中文乱码
		f.WriteString("\xEF\xBB\xBF")
	}
	w := csv.NewWriter(f)
	t := reflect.TypeOf(data).Elem().Elem()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	null_num := 0
	var map_arr = make(map[int]map[int]string)
	csvHeard := make([]string, 0)
	for i := 0; i < t.NumField(); i++ {
		excel_name, ok := t.Field(i).Tag.Lookup("excel_name")
		if !ok {
			null_num++
			continue
		}
		enums_str := t.Field(i).Tag.Get("enums")
		if enums_str != "" {
			map_data := make(map[int]string)
			enums_arr := strings.Split(enums_str, ",")
			for _, v := range enums_arr {
				enums_v := strings.Split(v, ":")
				key, _ := strconv.Atoi(enums_v[0])
				map_data[key] = enums_v[1]
			}
			map_arr[i] = map_data
		}
		csvHeard = append(csvHeard, excel_name)
	}
	if alread == 0 {
		w.Write(csvHeard)
	}
	num := reflect.ValueOf(data).Elem().Len()
	if num > 0 {
		for i := 0; i < num; i++ {
			value := reflect.ValueOf(data).Elem().Index(i)
			if reflect.ValueOf(data).Elem().Index(i).Type().Kind() == reflect.Ptr {
				value = value.Elem()
			}
			t := value.Type()
			null_num := 0
			cellData := make([]string, 0)
			for x := 0; x < t.NumField(); x++ {
				_, ok := t.Field(x).Tag.Lookup("excel_name")
				if !ok {
					null_num++
					continue
				}
				cell_value := value.Field(x)
				enums_str := t.Field(x).Tag.Get("enums")
				if enums_str != "" {
					map_data := map_arr[x]
					if _, ok := map_data[int(value.Field(x).Int())]; ok {
						cell_value = reflect.ValueOf(map_data[int(value.Field(x).Int())])
					}
				}
				excel_time := t.Field(x).Tag.Get("excel_time")
				switch excel_time {
				case "time":
					val_str := fmt.Sprintf("%v", value.Field(x))
					cell_value = reflect.ValueOf(val_str[:19])
					break
				case "int", "int64":
					t := time.Unix(value.Field(x).Int(), 0).Format("2006-01-02 15:04:05")
					cell_value = reflect.ValueOf(t)
					break
				}
				cellData = append(cellData, fmt.Sprintf("%v", cell_value))
			}
			w.Write(cellData)
		}
	}
	w.Flush()
	return
}

func (e *Excel) MapSaveCsv(file string, data []map[string]interface{}) (err error) {
	if len(data) < 2 {
		return errors.New("数据有问题")
	}
	alread := 1
	f := &os.File{}
	if filetool.FileIsExist(file) {
		f, err = os.OpenFile(file, os.O_WRONLY|os.O_APPEND, 0666)
	} else {
		f, err = os.Create(file)
		if err != nil {
			return err
		}
		alread = 0
	}
	defer f.Close()
	// 写入UTF-8 BOM，防止中文乱码
	if alread == 0 {
		f.WriteString("\xEF\xBB\xBF")
	}
	w := csv.NewWriter(f)
	dataSort := data[0]
	csvHeard := make([]string, len(dataSort))
	for k, v := range data[1] {
		if sort, ok := dataSort[k]; ok {
			if cast.ToInt(sort) > len(csvHeard)-1 {
				return errors.New("表头错误")
			}
			csvHeard[cast.ToInt(sort)] = fmt.Sprintf("%v", v)
		}
	}
	if alread == 0 {
		w.Write(csvHeard)
	}
	if len(data) == 2 {
		w.Flush()
		return
	}
	tempData := data[2:]
	for _, temp := range tempData {
		cellData := make([]string, len(csvHeard))
		for k, v := range temp {
			if sort, ok := dataSort[k]; ok {
				cellData[cast.ToInt(sort)] = fmt.Sprintf("%v", v)
			}
		}
		w.Write(cellData)
	}
	w.Flush()
	return
}

func addStr(v string, add int32) string {
	x := []rune(v)
	xLen := len(x)
	if (x[xLen-1] + add) > 90 {
		if xLen == 1 {
			return addStr("AA", add-1)
		}
		if xLen == 2 {
			return addStr(string(x[0]+1)+"A", add-1)
		}
	}
	x[xLen-1] = x[xLen-1] + add
	return string(x)
}
