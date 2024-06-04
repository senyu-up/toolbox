package storage

import (
	"context"
	"fmt"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"strings"
)

type ImageAuditReq struct {
	AppKey     string // 游戏 应用 key
	DataId     string // 图片在业务方的Id
	Url        string // 图片地址
	BizType    string // 业务类型
	DetectType string // 检测类型
}

type ImageAuditResult struct {
	JobId  string `json:"job_id"`
	AppKey string `json:"app_key"` // 游戏 应用 key
	DataId string `json:"data_id"` // 图片标识，该字段在结果中返回原始内容，长度限制为512字节。
	Url    string `json:"url"`     // 图片地址

	Result int  `json:"result"` // 有效值：0（审核正常），1（判定为违规敏感文件），2（疑似敏感，建议人工复核）。
	Score  int  `json:"score"`  // 该字段为敏感文件的分值，有效值：0-100。 例如：色情 99，则表明该内容非常有可能属于色情内容。
	State  bool `json:"state"`  // 处理状态，true 则表示处理成功

	Label    string `json:"label"`     // 当 result 为 1或2 时，该字段为敏感文件的具体类型，有效值：porn（色情）、terrorist（暴恐）、politics（涉政）、ads（广告）、qrcode（二维码）、sensitive（敏感）。
	SubLabel string `json:"sub_label"` // 该图命中的二级标签结果。
	Category string `json:"category"`  // 该字段为 Label 的子集，表示审核命中的具体审核类别。
}

type TXImageAudit struct {
	conf   config.ImageAudit
	client *cos.Client
}

func NewTXImageAudit(conf *config.ImageAudit) *TXImageAudit {
	bu, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", conf.Bucket, conf.Regin))
	cu, _ := url.Parse(fmt.Sprintf("https://%s.ci.%s.myqcloud.com", conf.Bucket, conf.Regin))
	b := &cos.BaseURL{BucketURL: bu, ServiceURL: nil, BatchURL: nil, CIURL: cu, FetchURL: nil}
	return &TXImageAudit{
		client: cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  conf.SecretId,
				SecretKey: conf.SecretKey,
			},
		}),
		conf: *conf,
	}
}

// 同步单张审核图片
// doc: https://cloud.tencent.com/document/product/460/37318
// sdk: https://github.com/tencentyun/cos-go-sdk-v5
func (tia *TXImageAudit) AuditImage(ctx context.Context, req *ImageAuditReq) (*ImageAuditResult, error) {

	var opt2 = &cos.ImageRecognitionOptions{
		CIProcess:        "sensitive-content-recognition",
		DetectUrl:        req.Url,        // 图片地址
		DetectType:       req.DetectType, // 使用默认 审核策略
		BizType:          req.BizType,
		LargeImageDetect: 0,
		Async:            0, // 同步请求
		MaxFrames:        1,
		DataId:           req.AppKey + "_" + req.DataId, // 这里注意给到腾讯的 dataId 不能重复
	}
	if opt2.BizType == "" {
		opt2.BizType = tia.conf.BizType
	}
	re, resp, err := tia.client.CI.ImageAuditing(ctx, "", opt2)
	logger.Debug("tx image audit re %v resp %v err %v", re, resp, err)
	if err == nil && re != nil && re.State == "Success" {
		// 审核成功
		return &ImageAuditResult{
			JobId:  re.JobId,
			AppKey: req.AppKey,
			DataId: req.DataId,
			Url:    re.Url,

			Result: re.Result,
			Score:  re.Score,
			State:  re.State == "Success",

			Label:    re.Label,
			SubLabel: re.SubLabel,
			Category: re.Category,
		}, err
	} else {
		logger.Ctx(ctx).Error("tx image audit re %v resp %v err %v", *re, *resp, err)
		return nil, err
	}
}

// 同步多张审核，耗时操作，注意grpc超时时间
// doc: https://cloud.tencent.com/document/product/460/63594
// sdk: https://cloud.tencent.com/document/product/460/72954
func (tia *TXImageAudit) AuditImages(ctx context.Context, urls map[string]string, appKey string,
	bizType string, detectType string) ([]*ImageAuditResult, error) {

	var opt = &cos.BatchImageAuditingOptions{
		Input: []cos.ImageAuditingInputOptions{},
		Conf: &cos.ImageAuditingJobConf{
			BizType:    bizType,
			Async:      0,
			DetectType: detectType,
		},
	}
	if opt.Conf.BizType == "" {
		opt.Conf.BizType = bizType
	}
	for i, u := range urls {
		opt.Input = append(opt.Input, cos.ImageAuditingInputOptions{
			DataId: appKey + "_" + i,
			Url:    u,
		})
	}
	re, resp, err := tia.client.CI.BatchImageAuditing(ctx, opt)
	logger.Debug("tx images batch audit re %v resp %v err %v", *re, *resp, err)
	if err == nil && re != nil && re.JobsDetail != nil {
		// 审核成功
		var result = []*ImageAuditResult{}
		var reAppKey = appKey
		var dataId string
		for _, jd := range re.JobsDetail {
			if jd.DataId != "" {
				// 从返回的 dataId 中解析出 appKey 和 dataId
				if strs := strings.Split(jd.DataId, "_"); len(strs) >= 1 {
					reAppKey = strs[0]
					dataId = strs[1]
				}
			}
			var item = &ImageAuditResult{
				JobId:  jd.JobId,
				AppKey: reAppKey,
				DataId: dataId,
				Url:    jd.Url,

				Result: jd.Result,
				Score:  jd.Score,
				State:  jd.State == "Success",

				Label:    jd.Label,
				SubLabel: jd.SubLabel,
				Category: jd.Category,
			}
			result = append(result, item)
		}
		return result, err
	} else {
		logger.Ctx(ctx).Error("tx images audit re %v resp %v err %v", re, resp, err)
		return nil, err
	}
}
