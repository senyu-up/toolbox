package struct_tool

import (
	"errors"
	"github.com/iancoleman/strcase"
	"reflect"
	"strings"
	"sync"
)

type CacheStruct struct {
	name   string
	rv     *reflect.Value
	rt     *reflect.Type
	fields []*CacheField
}

func (c *CacheStruct) Type() *reflect.Type {
	return c.rt
}

func (c *CacheStruct) Value() *reflect.Value {
	return c.rv
}

func (c *CacheStruct) Name() string {
	return c.name
}

func (c *CacheStruct) Fields() []*CacheField {
	return c.fields
}

type StructCache struct {
	lock sync.Mutex
	m    map[reflect.Type]*CacheStruct
}

// NewStructCache
// @description New一个StructCache, 简称sc
func NewStructCache() *StructCache {
	return &StructCache{
		m: make(map[reflect.Type]*CacheStruct),
	}
}

type CacheField struct {
	Idx       int
	Name      string
	SnakeName string
	Type      *reflect.StructField
	Tags      map[string]FieldTagParam
}

func (s *StructCache) get(t reflect.Type) (data *CacheStruct, exists bool) {
	data, exists = s.m[t]

	return
}

func (s *StructCache) set(t reflect.Type, data *CacheStruct) {
	s.m[t] = data
}

var ErrNotStruct = errors.New("not a struct")

// ExtractStructCache
// @description 从缓存获取StructCache, 未命中则对结构体进行解析并缓存
func (s *StructCache) ExtractStructCache(current interface{}) (data *CacheStruct, err error) {
	rt := reflect.TypeOf(current)
	step := 0

	if rt.Kind() == reflect.Ptr {
		step++
		rt = rt.Elem()
	}

	if rt.Kind() != reflect.Struct {
		err = ErrNotStruct
		return
	}

	//  无锁的读
	data, exists := s.get(rt)
	if exists {
		return data, nil
	}
	// 加锁的写
	s.lock.Lock()
	defer s.lock.Unlock()

	rv := reflect.ValueOf(current)
	for i := 0; i < step; i++ {
		rv = rv.Elem()
	}
	// 重试获取数据
	data, exists = s.get(rt)
	if exists {
		return data, nil
	}

	data = &CacheStruct{}
	data.name = rt.Name()
	data.rv = &rv
	data.rt = &rt

	data.fields = make([]*CacheField, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tagStr := string(field.Tag)

		data.fields[i] = &CacheField{
			Idx:       i,
			Name:      field.Name,
			Type:      &field,
			SnakeName: strcase.ToSnake(field.Name),
			Tags:      travelTag(tagStr),
		}
	}
	s.set(rt, data)

	return
}

type FieldTagParam []*TagParam

// GetByIndex
// @description 基于索引序号获取tag的param值
func (f FieldTagParam) GetByIndex(i int) (value *TagParam, exists bool) {
	if i >= len(f) {
		return nil, false
	} else {
		return f[i], true
	}
}

// Get
// @description 通过param的key获取param值
func (f FieldTagParam) Get(key string) (value *TagParam, exists bool) {
	for i, _ := range f {
		if f[i].Key == key {
			return f[i], true
		}
	}

	return nil, false
}

type TagParam struct {
	Key   string
	Value string
}

func travelTag(tag string) map[string]FieldTagParam {
	if tag == "" {
		return nil
	}

	data := make(map[string]FieldTagParam, 2)
	// protobuf:"varint,1,opt,name=ChannelID,proto3" json:"ChannelID,omitempty"
	tagList := strings.Split(tag, " ")
	for i, _ := range tagList {
		if pos := strings.Index(tagList[i], ":"); pos >= 0 {
			key := tagList[i][:pos]
			data[key] = FieldTagParam{}

			paramPear := strings.Split(tagList[i][pos+1:len(tagList[i])-1], ",")
			for j, _ := range paramPear {
				if paramPear[j] == "" {
					continue
				}
				if pos2 := strings.Index(paramPear[j], "="); pos2 >= 0 {
					paramKey := strings.Trim(paramPear[j][:pos2], `"`)
					paramValue := strings.Trim(paramPear[j][pos2+1:], `"`)

					data[key] = append(data[key], &TagParam{
						Key:   paramKey,
						Value: paramValue,
					})
				} else {
					data[key] = append(data[key], &TagParam{
						Key: strings.Trim(paramPear[j], `"`),
					})
				}
			}
		}
	}

	return data
}
