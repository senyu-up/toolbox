package struct_tool

import (
	"bytes"
	"encoding/gob"
	"github.com/senyu-up/toolbox/tool/runtime"
	"github.com/senyu-up/toolbox/tool/su_slice"
	"reflect"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

//将src中的值拷贝给dst
//对象值拷贝

func DeepCopy(dst, src interface{}) error {
	if src == nil {
		return nil
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

// JsonStrToMap json string to map
func JsonStrToMap(data []byte) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	return m, jsoniter.Unmarshal(data, &m)
}

func MakeUpdateField(i interface{}) map[string]interface{} {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	field := make(map[string]interface{})
	for j := 0; j < t.NumField(); j++ {
		if !v.Field(j).IsZero() {
			field[snakeString(t.Field(j).Name)] = v.Field(j).Interface()
		}
	}
	return field
}

func snakeString(s string) string {
	char := make([]byte, 0, len(s)*2)
	num := len(s)
	for i := 0; i < num; i++ {
		if i > 0 {
			if !(isUp(s[i]) == isUp(s[i-1])) {
				if !isUp(s[i-1]) {
					char = append(char, '_')
				}
			} else if isUp(s[i]) && isUp(s[i-1]) && i < (num-2) {
				if !isUp(s[i+1]) {
					char = append(char, '_')
				}
			}
		}
		char = append(char, s[i])
	}
	return strings.ToLower(string(char))
}

func isUp(i byte) bool {
	return i >= 'A' && i <= 'Z'
}

// 判断结构体是否存在某个属性值
func IsExistField(i interface{}, field string) bool {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for j := 0; j < t.NumField(); j++ {
		if t.Field(j).Name == field {
			return true
		}
	}
	return false
}

// StructToMap in转化成map
func StructToMap(in interface{}) (map[string]interface{}, error) {
	data, err := jsoniter.Marshal(in)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	return m, jsoniter.Unmarshal(data, &m)
}

// MapToStruct map in 转换成 struct out
func MapToStruct(in interface{}, out interface{}) error {
	data, err := jsoniter.Marshal(in)
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(data, out)
}

// 获取结构体字段数据, 返回map类型数据
// data: 结构体值
// tag：”字段标签类型，如：json“
// ignoreFields: 需要忽略的字段
// notIgnoreZero: 如果值为空，亦然将其加入到map。 为true时，不加入
func StructToMapByTag(data interface{}, tag string, ignoreFields []string, notIgnoreZero bool) map[string]interface{} {
	ret := map[string]interface{}{}
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	if t.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ret
		}
		ot := reflect.ValueOf(data).Elem().Type()
		if ot.Kind() != reflect.Struct {
			return ret
		}
		t = t.Elem()
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		var name string
		if len(tag) > 0 {
			name = t.Field(i).Tag.Get(tag)
		} else {
			name = t.Field(i).Name
		}
		if len(name) == 0 || name == "-" {
			continue
		}
		var tagNames = strings.Split(name, ",")
		// 如果 tag 有逗号,第一个字符串
		if 0 < len(tagNames) {
			name = tagNames[0]
		}

		value := v.Field(i).Interface()
		if nil != ignoreFields && su_slice.InArray(name, ignoreFields) {
			// 强制过滤掉的字段
			continue
		} else if !notIgnoreZero && runtime.IsZero(value) {
			// 如果忽略0值，需要判断是0 不
			continue
		}
		ret[name] = value
	}
	return ret
}
