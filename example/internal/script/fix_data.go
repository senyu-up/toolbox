package script

import (
	"context"
	"github.com/senyu-up/toolbox/example/global"
	"github.com/senyu-up/toolbox/example/internal/model"
)

// FixData
//
//	@Description: 修复用户数据
//	@param ctx  body any true "-"
//	@param db  body any true "-"
//	@return aff
//	@return err
func FixData(ctx context.Context) (aff int64, err error) {
	var db = global.GetFacade().GetMysqlClient()
	var persons = []model.Person{}
	err = db.WithContext(ctx).Model(model.Person{}).Where("age <= ?", 18).Find(&persons).Error
	if err != nil {
		return 0, err
	}
	for _, p := range persons {
		if p.Email == "" {
			err = db.WithContext(ctx).Model(model.Person{}).
				Where("id = ?", p.ID).UpdateColumn("email", "xx").Error
			if err != nil {
				return
			} else {
				aff++
			}
		}
	}

	return aff, err
}
