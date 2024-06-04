package storage

import (
	"errors"
	"github.com/senyu-up/toolbox/tool/su_slice"
	"io"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/logger"
)

type KeyFormat func(now time.Time, appKey, fileName string) string

type S3Conn struct {
	conf map[string]config.S3Storage

	sessions map[string]*session.Session

	//
	//  defaultRegin
	//  @Description: 默认 regin， 如果你指定了一个regin，那么就会使用这个regin，当你初始化了多个region时，请在调用接口时指定
	//
	defaultRegin string

	initRegions []string

	keyFormatter KeyFormat // key 格式化函数

	session *session.Session
	hkSess  *session.Session
	uwSess  *session.Session
	cnSess  *session.Session
}

var _s3Cli *S3Conn
var s3Lock sync.RWMutex

// InitByConf
//
//	@Description: 通过配置文件初始化 aws s3
//	@param awsConfig  body any true "-"
func InitByConf(awsConfig *config.Aws) *S3Conn {
	_S3Cli := &S3Conn{conf: map[string]config.S3Storage{},
		sessions: map[string]*session.Session{}, initRegions: []string{}}
	sa := awsConfig.AwsAccessId
	ss := awsConfig.AwsAccessKey

	for _, s3Regin := range awsConfig.S3 {
		if su_slice.InArray(s3Regin.Region, []string{AWSRegion, AWSUwRegion, HKAWSRegion, CNAWSRegion}) {
			sess, err := session.NewSession(&aws.Config{
				Region:      aws.String(s3Regin.Region),
				Credentials: credentials.NewStaticCredentials(sa, ss, ""),
			})
			if err != nil {
				logger.Error("Initialize S3 session regin: %s failed!", s3Regin.Region)
				continue
			}
			logger.Info("Initialize S3 session regin: %s success!", s3Regin.Region)
			_S3Cli.sessions[s3Regin.Region] = sess
			_S3Cli.conf[s3Regin.Region] = s3Regin
			_S3Cli.defaultRegin = s3Regin.Region // 默认 regin
		}
	}
	s3Lock.Lock()
	defer s3Lock.Unlock()
	_s3Cli = _S3Cli
	return _s3Cli
}

// Get
//
//	@Description: 传入 region 字符，获取对应 client
//	@receiver p
//	@param region  body any true "-"
//	@return s
//	@return c
//	@return e
func (p *S3Conn) Get(region string) (s *session.Session, c config.S3Storage, e error) {
	s3Lock.RLock()
	defer s3Lock.RUnlock()
	if s, ok := p.sessions[region]; ok {
		if c, ok2 := p.conf[region]; ok2 {
			return s, c, nil
		}
	}
	return nil, config.S3Storage{}, errors.New("regin not found")
}

// UploadFileCDN
//
//	@Description: 上传文件到S3 cdn,
//	@receiver p
//	@param file  body any true "-"
//	@param fileName   body any true "指定文件名"
//	@param region   body any true "如果不传，则使用默认的 region"
//	@return key	返回保存到 cdn 的文件路径， 只是 path，不含 bucket名，不含 host
//	@return err
func (p *S3Conn) UploadFileCDN(file io.Reader, fileName string, region ...string) (key string, err error) {
	var selectRegion string = p.defaultRegin
	if region != nil && 0 < len(region) {
		selectRegion = region[0]
	}
	if sess, conf, err := p.Get(selectRegion); err == nil {
		uploader := s3manager.NewUploader(sess)
		var key string
		if p.keyFormatter == nil {
			key = defaultDdnPathFormat(time.Now(), conf.Path, fileName)
		} else {
			key = p.keyFormatter(time.Now(), conf.Path, fileName)
		}
		//key = "/" + conf.Path + "/" + nowDate + "/" + fileName
		_, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(conf.Bucket),
			Key:    aws.String(key),
			Body:   file,
		})
		logger.Debug("Call UploadFileCDN regin %s file_name: %s", selectRegion, fileName)
		if err != nil {
			return key, err
		}
		logger.Info("upload success!!!")
		return key, nil
	} else {
		return key, err
	}
}

// 获取客户访问 cdn 地址
func (p *S3Conn) GetCdnPath(cdnPath string, region ...string) string {
	var selectRegion string = p.defaultRegin
	if region != nil && 0 < len(region) {
		selectRegion = region[0]
	}
	if _, conf, err := p.Get(selectRegion); err == nil {
		return conf.Host + cdnPath
	}
	return ""
}

// 获取s3预加密文件url
// deprecated
func (p *S3Conn) GetPreSignUrlCdn(region string, fileName string, dateStr ...string) (string, string, error) {
	if sess, conf, err := p.Get(region); err == nil {
		var svc *s3.S3
		svc = s3.New(sess)
		var nowDate string
		if dateStr != nil {
			nowDate = dateStr[0]
		}
		filePath := "/" + conf.Path + "/" + nowDate + "/" + fileName
		req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(conf.Bucket),
			Key:    aws.String(filePath),
		})
		urlStr, err := req.Presign(time.Duration(conf.Expire) * time.Hour)
		if err != nil {
			logger.Error("Failed to sign request %v", err)
			return "", "", err
		}
		cdnUrl := conf.Host + filePath
		return urlStr, cdnUrl, err
	} else {
		logger.Error("GetPreSignUrlCdn by regin: "+region+" error: %v", err)
		return "", "", err
	}
}

// 传入s3path，获取s3预加密文件url
func (p *S3Conn) GetPreSignUrl(min time.Duration, cdnPath string, region ...string) (string, error) {
	var selectRegion string = p.defaultRegin
	if region != nil && 0 < len(region) {
		selectRegion = region[0]
	}
	if sess, conf, err := p.Get(selectRegion); err == nil {
		var svc *s3.S3
		svc = s3.New(sess)
		req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
			Bucket: aws.String(conf.Bucket),
			Key:    aws.String(cdnPath),
		})
		urlStr, err := req.Presign(min)
		if err != nil {
			logger.Error("Failed to sign request %v", err)
			return "", err
		}
		return urlStr, err
	} else {
		logger.Error("GetPreSignUrlCdn by regin: "+selectRegion+" error: %v", err)
		return "", err
	}
}

// SetKeyFormatter
//
//	@Description: 设置 cdn key 格式化函数
//	@receiver p
//	@param f  body any true "-"
func (p *S3Conn) SetKeyFormatter(f func(now time.Time, prePath, fileName string) string) {
	p.keyFormatter = f
}

// defaultDdnPathFormat
//
//	@Description: 默认的 cdn key 格式化函数
//	@param now  body any true "-"
//	@param prePath  body any true "-"
//	@param fileName  body any true "-"
//	@return string
func defaultDdnPathFormat(now time.Time, prePath, fileName string) string {
	var nowDate = now.Format("20060102")
	return "/" + prePath + "/" + nowDate + "/" + fileName
}
