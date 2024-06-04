package center

import (
	"github.com/gofiber/fiber/v2"
	toolEnum "github.com/senyu-up/toolbox/enum"
	"github.com/senyu-up/toolbox/example/internal/dao/center_service"
	"github.com/senyu-up/toolbox/example/internal/enum"
	"github.com/senyu-up/toolbox/example/internal/model"
	"github.com/senyu-up/toolbox/tool/http/fiber/controller"
	"github.com/senyu-up/toolbox/tool/logger"
)

var CenterCtl = &Center{}

type Center struct {
	controller.BaseController
}

// @tags center
// @summary  获取一个事件
// @accept application/json
// @success 200 {object} api.GetEventResp
// @router /center/app_dsns [GET]
func (t *Center) ListAppDsns(c *fiber.Ctx) (err error) {

	// 请求数据
	res, err := center_service.GetAppDsnsByCond(c.Context(), center_service.GetDB(), model.CenterServiceAppDsnsFilter{})
	if err != nil {
		logger.Error("%v", err)
		return t.ResponseJson(c, controller.JsonResponse{Code: toolEnum.InternalErrCode, Msg: enum.DbFailedErrorMsg + err.Error()})
	}

	//返回标准json数据x
	err = t.ResponseJson(c, controller.JsonResponse{
		Code: enum.SuccessCode,
		Data: res,
	})
	if err != nil {
		logger.Error("%v", err)
		return err
	}

	//若返回err，则报500错误，展示error.Error()
	return nil
}

// @tags center
// @summary  获取一个 appDsns
// @accept application/json
// @param Req body api.IdReq true "json入参"
// @success 200 {object} api.GetEventResp
// @router /center/app_dsns/:id [GET]
func (t *Center) GetAppDsns(c *fiber.Ctx) (err error) {
	id, err := c.ParamsInt("id", 0)
	if err != nil {
		logger.Error("%v", err)
		return t.ResponseJson(c, controller.JsonResponse{Code: toolEnum.ParamsErrCode, Msg: enum.ParamsErrMsg + err.Error()})
	}
	// 请求数据
	var filter = model.CenterServiceAppDsnsFilter{Id: int32(id)}
	res, err := center_service.GetAppDsnsOneByCond(c.Context(), center_service.GetDB(), filter)
	if err != nil {
		logger.Error("%v", err)
		return t.ResponseJson(c, controller.JsonResponse{Code: toolEnum.InternalErrCode, Msg: enum.DbFailedErrorMsg + err.Error()})
	}

	//返回标准json数据x
	err = t.ResponseJson(c, controller.JsonResponse{
		Code: enum.SuccessCode,
		Data: res,
	})
	if err != nil {
		logger.Error("%v", err)
		return err
	}

	//若返回err，则报500错误，展示error.Error()
	return nil
}
