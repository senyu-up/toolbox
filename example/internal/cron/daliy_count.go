package cron

import (
	"context"
	"github.com/senyu-up/toolbox/example/global"
	"github.com/senyu-up/toolbox/example/internal/model"
	"github.com/senyu-up/toolbox/tool/logger"
)

// 需要实现 cron.Job 接口
func DailyCountJob(ctx context.Context) error {
	var err error
	var db = global.GetFacade().GetMysqlClient()
	var persons = []model.Person{}
	err = db.WithContext(ctx).Model(model.Person{}).Where("age <= ?", 18).Find(&persons).Error
	if err != nil {
		logger.Error("DailyCountJob Run, query person error", err)
		return err
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
	return nil
}
