package convert

import (
	"strconv"
)

func StringToInt32(v string) (int32, error) {
	i, err := strconv.ParseInt(v, 10, 32)
	return int32(i), err
}
