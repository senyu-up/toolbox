package page

import (
	"context"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type Page struct {
	conf  *config.MysqlPageConfig
	pg    int
	ps    int
	order string
	where []*Where
	param map[string]interface{}
	// 表名
	table string
	// 提取的字段
	fields   []string
	tplWhere string
	rawWhere string
}

const DftMaxPageSize int = 200
const DftPrimaryKey = "id"
const DftTotalKey = "total"

func NewConf(db ...*gorm.DB) *config.MysqlPageConfig {
	c := &config.MysqlPageConfig{
		MaxPageSize: DftMaxPageSize,
		PrimaryKey:  DftPrimaryKey,
		TotalKey:    DftTotalKey,
		DB:          nil,
	}

	if db != nil {
		c.DB = db[0]
	}

	return c
}

func New(conf *config.MysqlPageConfig) *Page {
	if conf.MaxPageSize == 0 {
		conf.MaxPageSize = DftMaxPageSize
	}
	if conf.PrimaryKey == "" {
		conf.PrimaryKey = DftPrimaryKey
	}
	if conf.TotalKey == "" {
		conf.TotalKey = DftTotalKey
	}

	return &Page{
		conf: conf,
		pg:   1,
		ps:   20,
	}
}

func (p *Page) Table(name string) *Page {
	p.table = name

	return p
}

func (p *Page) Fields(fields []string) *Page {
	p.fields = fields

	return p
}

func (p *Page) OrderBy(o string) *Page {
	p.order = o

	return p
}

func (p *Page) Where(cond []*Where) *Page {
	p.where = append(p.where, cond...)

	return p
}

func (p *Page) GetPage() int {
	return p.pg
}

func (p *Page) GetPageSize() int {
	return p.ps
}

type Where struct {
	// 当前的key
	Key string
	// 可选, 默认值, 如果不填, 当检测到当前 key 未在 param 变量中定义时, 会忽略本条件
	Default interface{}
	// 可选, 当定义了key但未定义sql时, sql会自动创建为 `key = :key` 形式
	Sql string
	// 可选, 用于判断当前的值条件是有有效, 如 nil, [] 都认为是无效值
	Ignore func(v interface{}, param map[string]interface{}) bool
}

func isZero(v interface{}, param map[string]interface{}) (b bool) {
	defer func() {
		if err := recover(); err != nil {
			b = true
		}
	}()
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)

	if rv.IsZero() {
		return true
	}
	// 对切片和数组进行长度判断
	if l, ok := v.([]interface{}); ok {
		return len(l) == 0
	}

	return false
}

func (p *Page) Limit(pg int, ps int) *Page {
	p.pg = pg
	p.ps = ps
	if p.pg < 1 {
		p.pg = -1
	} else if p.conf.MaxPage > 0 && p.pg > p.conf.MaxPage {
		p.pg = p.conf.MaxPage
	}

	if p.ps < 1 {
		p.ps = 20
	} else if p.conf.MaxPageSize > 0 && p.ps > p.conf.MaxPageSize {
		p.ps = p.conf.MaxPageSize
	}

	return p
}

// BuildWhere
// @description 构建where模版
func (p *Page) BuildWhere(param map[string]interface{}) string {
	w := strings.Builder{}
	for i, _ := range p.where {
		cKey := p.where[i].Key
		if cKey != "" {
			if curV, exists := param[cKey]; !exists {
				if p.where[i].Default != nil {
					param[cKey] = p.where[i].Default
				} else {
					// 未定义则忽略
					continue
				}
			} else {
				ignore := isZero
				if p.where[i].Ignore != nil {
					ignore = p.where[i].Ignore
				}

				if ignore(curV, param) {
					continue
				}
			}
			if p.where[i].Sql == "" {
				p.where[i].Sql = cKey + "=:" + cKey
			}
		} else {
			if p.where[i].Sql == "" {
				continue
			}
		}

		if w.Len() > 0 {
			w.WriteString(" AND ")
			w.WriteString(p.where[i].Sql)
		} else {
			w.WriteString(p.where[i].Sql)
		}
	}

	srcWhere := w.String()

	return srcWhere
}

func (p *Page) DoBuildWhere(param map[string]interface{}) (srcWhere string, where string, valList []interface{}, err error) {
	srcWhere = p.BuildWhere(param)
	where, valList, err = WhereParser(srcWhere, param)
	p.rawWhere = where
	p.tplWhere = srcWhere

	return
}

func (p *Page) List(ctx context.Context, param map[string]interface{}) (rs *Result) {
	tx := p.conf.DB.WithContext(ctx).Table(p.table)
	if param == nil {
		param = make(map[string]interface{})
	}
	p.param = param
	if p.fields != nil {
		tx.Select(p.fields)
	}
	if p.pg > 0 {
		offset := (p.pg - 1) * p.ps
		tx.Offset(offset)
	}

	tx.Limit(p.ps)
	tx.Order(p.order)
	srcWhere, where, valList, err := p.DoBuildWhere(p.param)
	if srcWhere != "" {
		tx.Where(where, valList...)
	}
	rs = &Result{
		err:      err,
		p:        p,
		rawWhere: where,
		tplWhere: srcWhere,
		valList:  valList,
		tx:       tx,
	}

	return rs
}

const variablePrefix = ':'

func WhereParser(sql string, param map[string]interface{}) (where string, valList []interface{}, err error) {
	from := strings.IndexByte(sql, variablePrefix)
	if from == -1 {
		return sql, valList, nil
	}
	valList = make([]interface{}, 0, len(param))

	l := len(sql)

	var v strings.Builder
	var vFlag bool
	var w strings.Builder

	w.WriteString(sql[:from])
	for ; from < l; from++ {
		if sql[from] == variablePrefix {
			v.Reset()
			vFlag = true
			w.WriteByte('?')
		} else {
			if vFlag {
				i := int(sql[from])
				if (i >= 65 && i <= 90) || (i >= 97 && i <= 122) || (i >= 48 && i <= 57) || (i == 95) {
					v.WriteByte(sql[from])
				} else {
					w.WriteByte(sql[from])
					k := v.String()
					vFlag = false
					if k == "" {
						continue
					}
					v.Reset()
					// 处理变量
					if pv, exists := param[k]; exists {
						valList = append(valList, pv)
					} else {
						err = errors.New("param key not defined key:[" + k + "]")
						return
					}
				}
			} else {
				w.WriteByte(sql[from])
			}
		}
	}

	if v.Len() > 0 {
		// 处理尾巴变量
		k := v.String()
		// 处理变量
		if pv, exists := param[k]; exists {
			valList = append(valList, pv)
		} else {
			err = errors.New("param key not defined key:" + k)
			return
		}
	}

	return w.String(), valList, nil
}

type Result struct {
	p        *Page
	err      error
	tplWhere string
	rawWhere string
	valList  []interface{}
	tx       *gorm.DB
}

func (r *Result) GetPage() int {
	return r.p.pg
}

func (r *Result) GetPageSize() int {
	return r.p.ps
}

func (r *Result) Error() error {
	if r.err != nil {
		return r.err
	}
	return r.tx.Error
}

func (r *Result) Scan(data interface{}) error {
	if r.err != nil {
		return r.err
	}
	if r.tx == nil {
		return errors.New("nil transaction")
	}

	return r.tx.Scan(data).Error
}

func (r *Result) Count(f ...func(where string) string) (cnt int64, err error) {
	var cntRs []int64
	if f != nil {
		var whereSql = r.tplWhere
		if whereSql == "" {
			whereSql = "1"
		}
		sql, valList, err := WhereParser(f[0](whereSql), r.p.param)
		if err != nil {
			return 0, err
		}
		err = r.tx.Raw(sql, valList...).Pluck(r.p.conf.TotalKey, &cntRs).Error
		if cntRs != nil {
			cnt = cntRs[0]
		}

		return cnt, err
	} else {
		where := r.rawWhere
		if where == "" {
			where = "1"
		}
		err = r.tx.Raw("SELECT count(0) AS "+r.p.conf.TotalKey+" FROM "+r.p.table+" WHERE "+where, r.valList...).Pluck(r.p.conf.TotalKey, &cntRs).Error
		if len(cntRs) > 0 {
			cnt = cntRs[0]
		}

		return cnt, err
	}
}

func (p *Page) Raw(ctx context.Context, rawSql string, valList []interface{}) (rs *Result) {
	tx := p.conf.DB.WithContext(ctx).Raw(rawSql, valList...)
	rs = &Result{
		p:        p,
		tplWhere: p.tplWhere,
		rawWhere: p.rawWhere,
		err:      tx.Error,
		tx:       tx,
	}

	// set empty
	p.tplWhere = ""
	p.rawWhere = ""

	return rs
}

func (p *Page) Query(ctx context.Context, sql string, param map[string]interface{}) (rs *Result) {
	p.param = param
	sqlTpl, valList, err := WhereParser(sql, param)
	rs = &Result{
		p:        p,
		err:      err,
		tplWhere: sql,
		rawWhere: sqlTpl,
		valList:  valList,
	}
	if err != nil {
		return
	}

	rs.tx = p.conf.DB.WithContext(ctx).Raw(sqlTpl, valList...)

	return rs
}

func ToParamData(param interface{}) (data map[string]interface{}) {
	byteData, _ := jsoniter.Marshal(param)
	data = make(map[string]interface{})
	gjson.ParseBytes(byteData).ForEach(func(key, value gjson.Result) bool {
		if value.Type == gjson.Number {
			// 这里可能是 int, 也可能是 float, 需要进一步处理
			if strings.Contains(value.String(), ".") {
				data[key.String()] = value.Float()
			} else {
				data[key.String()] = value.Int()
			}
		} else {
			data[key.String()] = value.Value()
		}
		return true
	})

	return data
}
