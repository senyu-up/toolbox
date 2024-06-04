package convert

import (
	"fmt"
	"strconv"
)

func UIntToString(v uint) string {
	return strconv.Itoa(int(v))
}

// Decimal
// @description 保留两位小数
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

// Decimal5
// @description 保留5位小数
func Decimal5(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.5f", value), 64)
	return value
}

// Decimal4
// @description 保留4位小数
func Decimal4(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", value), 64)
	return value
}

// Decimal32
// @description 保留两位小数
func Decimal32(num float32) float32 {
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", num), 64)
	return float32(value)
}

// Float2Float
// @description 保留8位小数
func Float2Float(num float64) float64 {
	floatNum, _ := strconv.ParseFloat(fmt.Sprintf("%.8f", num), 64)
	return floatNum
}
