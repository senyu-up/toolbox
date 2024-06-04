package su_slice

import (
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/cast"
)

// 判断Slice中是否含有某个元素
func IncludeItem(array interface{}, target interface{}) (int, bool) {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(target, s.Index(i).Interface()) == true {
				return i, true
			}
		}
	}
	return -1, false
}

// 查找字符是否在数组中
// @deprecated 推荐使用 NewMapping() 形式
func InArray(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}
func Implode(arr interface{}, separator string) string {
	inArr := cast.ToStringSlice(arr)
	return strings.Join(inArr, separator)
}

// Difference 取出差集
func Difference(slice1, slice2 []string) []string {
	m := make(map[string]string)
	for _, v := range slice1 {
		m[v] = v
	}
	for _, v := range slice2 {
		if m[v] != "" {
			delete(m, v)
		}
	}
	var str []string
	for _, s2 := range m {
		str = append(str, s2)
	}
	return str
}

type SliceNRs [][]interface{}

func (s SliceNRs) StringSlice() [][]string {
	rs := make([][]string, 0, len(s))
	for i, _ := range s {
		var s1 = make([]string, 0, len(s[i]))
		for j, _ := range s[i] {
			s1 = append(s1, cast.ToString(s[i][j]))
		}

		rs = append(rs, s1)
	}

	return rs
}

func (s SliceNRs) IntSlice() [][]int {
	rs := make([][]int, 0, len(s))
	for i, _ := range s {
		var s1 = make([]int, 0, len(s[i]))
		for j, _ := range s[i] {
			s1 = append(s1, cast.ToInt(s[i][j]))
		}

		rs = append(rs, s1)
	}

	return rs
}

func (s SliceNRs) Int32Slice() [][]int32 {
	rs := make([][]int32, 0, len(s))
	for i, _ := range s {
		var s1 = make([]int32, 0, len(s[i]))
		for j, _ := range s[i] {
			s1 = append(s1, cast.ToInt32(s[i][j]))
		}

		rs = append(rs, s1)
	}

	return rs
}

func (s SliceNRs) Int64Slice() [][]int64 {
	rs := make([][]int64, 0, len(s))
	for i, _ := range s {
		var s1 = make([]int64, 0, len(s[i]))
		for j, _ := range s[i] {
			s1 = append(s1, cast.ToInt64(s[i][j]))
		}

		rs = append(rs, s1)
	}

	return rs
}

func (s SliceNRs) Float32Slice() [][]float32 {
	rs := make([][]float32, 0, len(s))
	for i, _ := range s {
		var s1 = make([]float32, 0, len(s[i]))
		for j, _ := range s[i] {
			s1 = append(s1, cast.ToFloat32(s[i][j]))
		}

		rs = append(rs, s1)
	}

	return rs
}

func (s SliceNRs) Float64Slice() [][]float64 {
	rs := make([][]float64, 0, len(s))
	for i, _ := range s {
		var s1 = make([]float64, 0, len(s[i]))
		for j, _ := range s[i] {
			s1 = append(s1, cast.ToFloat64(s[i][j]))
		}

		rs = append(rs, s1)
	}

	return rs
}

// SplitN 将一个大切片拆成多个小切片
func SplitN(slice interface{}, size int) SliceNRs {
	rv := reflect.ValueOf(slice)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		var l = rv.Len()
		var sliceNum = l/size + 1
		rs := make(SliceNRs, 0, sliceNum)
		for i := 0; i < l; i += size {
			curS := make([]interface{}, 0, size)
			for j := 0; j < size; j++ {
				if i+j >= l {
					break
				}
				curS = append(curS, rv.Index(i+j).Interface())
			}

			rs = append(rs, curS)
		}

		return rs
	}

	return nil
}

// Intersection 取出交集
func Intersection(a []string, b []string) (inter []string) {
	m := make(map[string]string)
	nn := make([]string, 0)
	for _, v := range a {
		m[v] = v
	}
	for _, v := range b {
		times, _ := m[v]
		if len(times) > 0 {
			nn = append(nn, v)
		}
	}
	return nn
}

func RemoveDuplicate(list *[]int32) []int32 {
	x := []int32{}
	for _, i := range *list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}

type SliceIndex map[interface{}]int

func Index(v interface{}) SliceIndex {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		m := make(SliceIndex, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			m[rv.Index(i).Interface()] = i
		}

		return m
	}

	return nil
}

/*Pos
* @Description: 获取key在切片中的位置
* @param key
* @return int
 */
func (m SliceIndex) Pos(key interface{}) int {
	if m == nil {
		return -1
	}
	pos, exists := m[key]
	if !exists {
		pos = -1
	}
	return pos
}

/*Has
* @Description: 判断key是否存在
* @param key
* @return bool
 */
func (m SliceIndex) Has(key interface{}) bool {
	if m == nil {
		return false
	}
	_, exists := m[key]

	return exists
}

type SliceStructMapping map[interface{}]interface{}

func IndexSliceStruct(v interface{}, field string) SliceStructMapping {
	var m SliceStructMapping
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	var ptrFlag bool
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		m = make(SliceStructMapping, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			curItem := rv.Index(i)
			if i == 0 {
				if curItem.Kind() == reflect.Ptr {
					ptrFlag = true
				}
			}
			if !curItem.IsZero() {
				var curField reflect.Value
				if ptrFlag {
					curField = curItem.Elem().FieldByName(field)
				} else {
					curField = curItem.FieldByName(field)
				}
				if curField.CanAddr() {
					m[curField.Interface()] = curItem.Interface()
				}
			}
		}

		return m
	}

	return nil
}

// Get
// @description 基于索引的field名进行查询, 返回对应的数据, 如果查询不到  ok=false
func (s SliceStructMapping) Get(k interface{}) (data interface{}, ok bool) {
	if s == nil {
		return nil, false
	}
	if v, ok := s[k]; ok {
		return v, true
	} else {
		return nil, false
	}
}

// Has
// @description 判断当前的k是否存在
func (s SliceStructMapping) Has(k interface{}) bool {
	if s == nil {
		return false
	}

	_, ok := s[k]

	return ok
}

var Rander = rand.New(rand.NewSource(time.Now().UnixNano()))

// 获取一定范围(min,max) 的随机数
func GetRandomNum(min, max int) int {
	return Rander.Intn(max-min) + min
}
