package db

import (
	"context"
	"github.com/senyu-up/toolbox/tool/logger"
	"github.com/senyu-up/toolbox/tool/trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"log"
	"sync"
	"time"
)

type MongoMonitor struct {
	logLevel      int           // 5=> debug, 4=> info, 3=> warn, 2=> error, 1=> silent
	traceOn       bool          // 开启 链路 ？
	slowThreshold time.Duration // 慢日志阈值
	split         int           // 分多少层级
	splitDuration time.Duration // 每层级的时间间隔

	monitorLogMap sync.Map // reqId -> event
	seqLogId      sync.Map // m log id -> time unix

	useLogCount uint32 // 使用的日志次数
	gcCount     uint32 // gc 频次
}

func (m *MongoMonitor) addLogId(ctx context.Context, defaultLog logger.Log, evt *event.CommandStartedEvent) {
	if es, ok2 := m.monitorLogMap.Load(evt.RequestID); ok2 {
		defaultLog.Ctx(ctx).Warn("mongo monitor duplicate reqId:%d cmd:%s dEvent: %+v",
			evt.RequestID, evt.CommandName, es)
	}
	m.monitorLogMap.Store(evt.RequestID, evt)
	m.seqLogId.Store(evt.RequestID, time.Now().Unix())
}

func (m *MongoMonitor) popLogById(ctx context.Context, defaultLog logger.Log, reqId int64) (evt *event.CommandStartedEvent) {
	defer func() {
		m.useLogCount++
		if m.gcCount > 0 && m.useLogCount%m.gcCount == 0 {
			var gcTime = time.Now().Add(-time.Hour).Unix() // 倒退1小时
			m.seqLogId.Range(func(key, value interface{}) bool {
				if value.(int64) < gcTime { // 如果日志时间早于gc时间，则删除
					m.seqLogId.Delete(key)
					m.monitorLogMap.Delete(key)
				}
				return true
			})
		}
	}()
	if es, ok2 := m.monitorLogMap.Load(reqId); !ok2 {
		defaultLog.Ctx(ctx).Warn("mongo monitor not found reqId:%d", reqId)
		return nil
	} else {
		m.monitorLogMap.Delete(reqId)
		m.seqLogId.Delete(reqId)
		return es.(*event.CommandStartedEvent)
	}
}

func (m *MongoMonitor) GetMonitor(defaultLog logger.Log) *event.CommandMonitor {
	return &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			m.addLogId(ctx, defaultLog, evt)
			if m.logLevel > 4 {
				if defaultLog != nil {
					defaultLog.Ctx(ctx).Debug("mongo db:%s reqId:%d cmd:%s exec: %v", evt.DatabaseName,
						evt.RequestID, evt.CommandName, evt.Command)
				} else {
					log.Printf("mongo db:%s reqId:%d cmd:%s exec: %v", evt.DatabaseName,
						evt.RequestID, evt.CommandName, evt.Command)
				}
			}
		},
		Succeeded: func(ctx context.Context, evt *event.CommandSucceededEvent) {
			var cmdRaw bson.Raw
			var execDur = time.Duration(evt.DurationNanos) * time.Nanosecond
			if re := m.popLogById(ctx, defaultLog, evt.RequestID); re != nil {
				cmdRaw = re.Command
			}
			// 开启链路追踪, 且 traceId 不为空，才记录
			if m.traceOn && len(trace.ParseRequestId(ctx)) > 0 {
				var traceId, spanId, newSpanId string
				ctx, traceId, spanId, newSpanId = trace.ParseOrGenContext(ctx)
				span := trace.NewJaegerSpan("mongo", traceId, newSpanId, spanId, nil, nil)
				if cmdRaw != nil {
					// 找到了 sql 记录
					span.SetTag("bson.raw", cmdRaw)
				}
				if execDur > m.slowThreshold {
					span.SetTag("slowsql", 1)
				}
				span.Finish()
			}
			if execDur > m.slowThreshold {
				// 慢日志 告警
				if 2 < m.logLevel {
					if defaultLog != nil {
						defaultLog.Ctx(ctx).Warn("mongo reqId:%d cmd:%s dur:%dms sql:%s", evt.RequestID, evt.CommandName,
							execDur.Milliseconds(), cmdRaw)
					} else {
						log.Printf("mongo reqId:%d cmd:%s dur:%dms sql:%s", evt.RequestID, evt.CommandName,
							execDur.Milliseconds(), cmdRaw)
					}
				}
			} else {
				if 3 < m.logLevel {
					if defaultLog != nil {
						defaultLog.Ctx(ctx).Info("mongo reqId:%d cmd:%s dur:%dms sql:%s", evt.RequestID, evt.CommandName,
							execDur.Milliseconds(), cmdRaw)
					} else {
						log.Printf("mongo reqId:%d cmd:%s dur:%dms sql:%s", evt.RequestID, evt.CommandName,
							execDur.Milliseconds(), cmdRaw)
					}
				}
			}
		},
		Failed: func(ctx context.Context, evt *event.CommandFailedEvent) {
			var cmdRaw bson.Raw
			var execDur = time.Duration(evt.DurationNanos) * time.Nanosecond
			if re := m.popLogById(ctx, defaultLog, evt.RequestID); re != nil {
				cmdRaw = re.Command
			}
			// 开启链路追踪, 且 traceId 不为空，才记录
			if m.traceOn && len(trace.ParseRequestId(ctx)) > 0 {
				var traceId, spanId, newSpanId string
				ctx, traceId, spanId, newSpanId = trace.ParseOrGenContext(ctx)
				span := trace.NewJaegerSpan("mongo", traceId, newSpanId, spanId, nil, nil)
				span.SetTag("error", evt.Failure)
				if cmdRaw != nil {
					// 找到了 sql 记录
					span.SetTag("bson.raw", cmdRaw)
				}
				if execDur > m.slowThreshold {
					span.SetTag("slowsql", 1)
				}
				span.Finish()
			}
			if 1 < m.logLevel {
				if defaultLog != nil {
					defaultLog.Ctx(ctx).Error("mongo reqId:%d cmd:%s dur:%dms failure:%s sql:%s", evt.RequestID,
						evt.CommandName, execDur.Milliseconds(), evt.Failure, cmdRaw)
				} else {
					log.Printf("mongo reqId:%d cmd:%s dur:%dms failure:%s sql:%s", evt.RequestID,
						evt.CommandName, execDur.Milliseconds(), evt.Failure, cmdRaw)
				}
			}
		},
	}
}

// GetDefaultMonitor 获取简单 默认 mongo log monitor
func GetDefaultMonitor(defaultLog logger.Log) *event.CommandMonitor {
	return &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			if defaultLog != nil {
				defaultLog.Ctx(ctx).Info("mongo db:%s reqId:%d cmd:%s exec: %v", evt.DatabaseName,
					evt.RequestID, evt.CommandName, evt.Command)
			} else {
				log.Printf("mongo db:%s reqId:%d cmd:%s exec: %v", evt.DatabaseName,
					evt.RequestID, evt.CommandName, evt.Command)
			}
		},
		Succeeded: func(ctx context.Context, evt *event.CommandSucceededEvent) {
			var execDur = time.Duration(evt.DurationNanos) * time.Nanosecond
			if defaultLog != nil {
				defaultLog.Ctx(ctx).Error("mongo reqId:%d cmd:%s dur:%d ms", evt.RequestID, evt.CommandName,
					execDur.Milliseconds())
			} else {
				log.Printf("mongo reqId:%d cmd:%s dur:%d ms", evt.RequestID, evt.CommandName,
					execDur.Milliseconds())
			}
		},
		Failed: func(ctx context.Context, evt *event.CommandFailedEvent) {
			var execDur = time.Duration(evt.DurationNanos) * time.Nanosecond
			if defaultLog != nil {
				defaultLog.Ctx(ctx).Error("mongo reqId:%d cmd:%s dur:%d ms failure:%s", evt.RequestID,
					evt.CommandName, execDur.Milliseconds(), evt.Failure)
			} else {
				log.Printf("mongo reqId:%d cmd:%s dur:%d ms failure:%s", evt.RequestID,
					evt.CommandName, execDur.Milliseconds(), evt.Failure)
			}
		},
	}
}
