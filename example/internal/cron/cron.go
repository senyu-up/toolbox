package cron

import (
	"context"
	"github.com/senyu-up/toolbox/example/global"
	"github.com/senyu-up/toolbox/example/internal/model"
	"github.com/senyu-up/toolbox/tool/runtime"
	"time"
)

// 用户日常统计
func UserDalyStatistics(ctx context.Context) (err error) {
	runtime.GoWithPanic(ctx, "user_daly_statistics", func() {
		for {
			time.Sleep(time.Minute)
			var db = global.GetFacade().GetMysqlClient()
			var persons = []model.Person{}
			err = db.WithContext(ctx).Model(model.Person{}).Where("age <= ?", 18).Find(&persons).Error
			if err != nil {
				continue
			}
			for _, p := range persons {
				if p.Email == "" {
					err = db.WithContext(ctx).Model(model.Person{}).
						Where("id = ?", p.ID).UpdateColumn("email", "xx").Error
					if err != nil {
						continue
					}
				}
			}

			time.Sleep(time.Hour * 24) // 一天执行一次
		}
	})
	return err
}
