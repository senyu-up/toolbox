package event

import "errors"

var (
	//DoesNotExistErr 事件不存在
	DoesNotExistErr = errors.New("event does not exist")
)
