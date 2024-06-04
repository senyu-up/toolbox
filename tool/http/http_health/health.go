package http_health

import (
	"fmt"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/logger"
	"net/http"
	"net/http/pprof"
	"time"
)

type HealthChecker struct {
	conf   config.HealthCheck
	server http.Server
}

func PprofHandle(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

func ok(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(("ok")))
}

func okPrint(w http.ResponseWriter, r *http.Request) {
	logger.Info(fmt.Sprintf("%s: %d", "system health:", time.Now().Unix()))
	w.Write([]byte(("ok")))
}

// 如果你的项目内没有使用fiber, 但又需要上报 pprof 可以使用这个函数
func NewHttpHealthCheckServer(opts ...HealthOption) (*HealthChecker, error) {
	var conf = config.HealthCheck{
		Addr:  "0.0.0.0",
		Pprof: true,
		Port:  80,
	}
	for _, opt := range opts {
		opt(&conf)
	}

	addr := conf.Addr
	if conf.Port == 0 {
		return nil, fmt.Errorf("health checker port is invalid")
	}
	addr += ":" + fmt.Sprintf("%d", conf.Port)
	server := http.Server{
		Addr:        addr,
		ReadTimeout: 6 * time.Second,
	}
	mux := http.NewServeMux()
	if conf.DisableLog {
		mux.HandleFunc("/system/health", ok)
	} else {
		mux.HandleFunc("/system/health", okPrint)
	}
	server.Handler = mux
	if conf.Pprof {
		PprofHandle(mux)
		logger.Info("http health checker pprof enabled")
	} else {
		logger.Info("http health checker pprof disabled")
	}

	return &HealthChecker{
		conf:   conf,
		server: server,
	}, nil
}

// Start 启动健康检查， 阻塞运行
func (h *HealthChecker) Start() error {
	return h.server.ListenAndServe()
}

func (h *HealthChecker) Addr() string {
	return h.conf.Addr
}

func (h *HealthChecker) Close() error {
	return h.server.Close()
}

type HandleFunc func() bool

// HttpHandle http健康检测
func HttpHandle(funs ...HandleFunc) {
	http.HandleFunc("/system/health", func(w http.ResponseWriter, r *http.Request) {
		for _, f := range funs {
			if !f() {
				logger.Warn(fmt.Sprintf("%s: %d", "system health close:", time.Now().Unix()))
				w.WriteHeader(500)
				return
			}
		}
		logger.Warn(fmt.Sprintf("%s: %d", "system health:", time.Now().Unix()))
		_, _ = w.Write([]byte(("ok")))
	})
}
