package script

import (
	"context"
	"fmt"
	"github.com/senyu-up/toolbox/example/global"
)

func PushEvent(args map[string]string) error {
	fmt.Printf("get args %v\n", args)
	var topic = "xh_push"
	var ctx = context.TODO()
	p, o, err := global.GetFacade().GetKafkaClient().PushSync(ctx, topic, args)
	fmt.Printf("push result %v %v %v\n", p, o, err)
	return nil
}
