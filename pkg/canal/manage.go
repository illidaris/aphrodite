package canal

import (
	"context"
	"sync"

	"github.com/illidaris/aphrodite/pkg/contextex"
)

var globLog ILogger = DefaultLogger{}

func SetLogger(log ILogger) {
	globLog = log
}

var once sync.Once
var manager *ConnectManager

func NewConnectManager() *ConnectManager {
	once.Do(func() {
		manager = &ConnectManager{
			Connectors:   []ISyncConnector{},
			ConnectorMap: map[string]ISyncConnector{},
		}
	})
	return manager
}

type ConnectManager struct {
	sync.RWMutex
	Connectors   []ISyncConnector
	ConnectorMap map[string]ISyncConnector
}

func (i *ConnectManager) AddConnector(raw context.Context, connector ISyncConnector) {
	i.Lock()
	defer i.Unlock()
	if i.ConnectorMap == nil {
		i.ConnectorMap = make(map[string]ISyncConnector)
	}
	id := connector.Id()
	i.Connectors = append(i.Connectors, connector)
	i.ConnectorMap[id] = connector
	globLog.Info(raw, "ConnectManager_AddConnector %s", id)
	go func(subCtx context.Context, v string) {
		subCtx = contextex.TransferBackground(subCtx)
		defer func() {
			if r := recover(); r != nil {
				globLog.Error(subCtx, "[PANIC] %s %v", v, r)
			}
		}()
		globLog.Info(subCtx, "ConnectManager_AddConnector %s Go", id)
		err := connector.Run(subCtx)
		if err != nil {
			globLog.Error(subCtx, "Connect_Run_Err %s %v", v, err)
			return
		}
		defer func() {
			connector.Close(subCtx)
		}()
	}(raw, id)
}

func (i *ConnectManager) Stats(ctx context.Context) map[string]OperateStats {
	i.RLock()
	defer i.RUnlock()
	stats := map[string]OperateStats{}
	for k, v := range i.ConnectorMap {
		stats[k] = v.Stats(ctx)
	}
	return stats
}

func (i *ConnectManager) Stat(ctx context.Context, id string) OperateStats {
	i.RLock()
	defer i.RUnlock()
	v, ok := i.ConnectorMap[id]
	if !ok {
		return OperateStats{}
	}
	return v.Stats(ctx)
}
