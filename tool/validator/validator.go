package validator

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
)

var (
	IPRegexp         = regexp.MustCompile(`^(([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.)(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){2}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	xxxVersionRegexp = regexp.MustCompile(`^([1-9]\d|[1-9])(\.([1-9]\d|\d))(\.1[\d]{9})$`)
	PhoneRegexp      = regexp.MustCompile(`^(1[3-9][0-9]|14[5|7]|15[0|1|2|3|5|6|7|8|9]|18[0|1|2|3|5|6|7|8|9]|177)\d{8}$`)
	MailRegex        = regexp.MustCompile(`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`)
)

// AppKey校验规则
// ZGV2ZWxvcDpDSmpCNnpEVWFhOjE2NDMzNDQ0MDI=
// develop:CJjB6zDUaa:1643344402
// base64decode
var appKeyPattern = regexp.MustCompile(`^[a-zA-Z0-9]+:[a-zA-Z0-9]+:\d{10}$`)

// 订单ID校验规则
// 生成方式："XH_" + xid.New().String()
// 示例：XH_c0vlnupfrd4ggl5t4aug
var xhOrderPattern = regexp.MustCompile(`XH_[a-zA-Z0-9]+`)

// NewSUValidator
// @description new一个validator, 带有 appKey、order 校验规则
func NewSUValidator() *validator.Validate {
	v := validator.New()
	_ = v.RegisterValidation("appKey", func(fl validator.FieldLevel) bool {

		s, err := base64.StdEncoding.DecodeString(fl.Field().String())
		if err != nil {
			return false
		}
		if appKeyPattern.Match(s) {
			return true
		}
		return false
	})
	_ = v.RegisterValidation("order", func(fl validator.FieldLevel) bool {
		if xhOrderPattern.MatchString(fl.Field().String()) {
			return true
		}
		return false
	})

	validate = v
	return v
}

func VerifyIp(ip string) error {
	match := IPRegexp.Match([]byte(ip))
	if match {
		return nil
	}

	return errors.New("ip format error")
}

func VerifyxxxVersion(version string) error {
	match := xxxVersionRegexp.Match([]byte(version))
	if match {
		return nil
	}

	return errors.New("xxx version format error")
}

func VerifyPhone(phone string) error {
	match := PhoneRegexp.Match([]byte(phone))
	if match {
		return nil
	}

	return errors.New("phone format error")
}

func VerifyMail(mail string) error {
	match := MailRegex.Match([]byte(mail))
	if match {
		return nil
	}

	return errors.New("mail format error")
}

var validate = validator.New()

// 检验结构体参数
func StructValidator(s interface{}) error {
	err := validate.Struct(s)

	return err
}

// IsNilInterface
// @description 判断是否为空interface{}
func IsNilInterface(data interface{}) bool {
	if data == nil {
		return true
	} else {
		rt := reflect.TypeOf(data)
		switch rt.Kind() {
		case reflect.Ptr, reflect.Map, reflect.Slice:
			return reflect.ValueOf(data).IsNil()
		default:
			return false
		}
	}
}

func IsNilValue(v interface{}) bool {
	if v == nil {
		return true
	} else {
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Ptr {
			if rv.Elem().Kind() == reflect.Ptr {
				rv = rv.Elem()
			}
		}

		return rv.IsZero()
	}
}

// EmptyError
// @description 判断data是否为空, 如果为空则按照 err, msg 的优先级返回对应的错误
func EmptyError(data interface{}, err error, module string, msg string) error {
	if err != nil {
		if module != "" {
			return fmt.Errorf("module:%s err:%s", module, err.Error())
		} else {
			return err
		}
	}

	if IsNilInterface(data) {
		if module != "" {
			return fmt.Errorf("moduel:%s err:%s", module, msg)
		} else {
			return errors.New(msg)
		}
	}

	return nil
}
