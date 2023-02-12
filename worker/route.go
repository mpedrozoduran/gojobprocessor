package worker

import (
	"fmt"
	"github.com/mpedrozoduran/gojobprocessor/proto/master"
	"sync"
)

type RouteTable struct {
	workerTable sync.Map
}

func (r *RouteTable) Push(info master.WorkerInfo) {
	r.workerTable.Store(info.Id, info)
}

func (r *RouteTable) Get(key string) (master.WorkerInfo, error) {
	if val, ok := r.workerTable.Load(key); ok {
		return val.(master.WorkerInfo), nil
	}
	return master.WorkerInfo{}, fmt.Errorf("worker with id=%s not found", key)
}

func (r *RouteTable) GetAll() map[string]master.WorkerInfo {
	res := make(map[string]master.WorkerInfo)
	r.workerTable.Range(func(key, value any) bool {
		res[key.(string)] = value.(master.WorkerInfo)
		return true
	})
	return res
}

func (r *RouteTable) Delete(key string) {
	r.workerTable.Delete(key)
}

func (r *RouteTable) Clear(keys map[string]string) {
	for k := range keys {
		r.workerTable.Delete(k)
	}
}

func (r *RouteTable) Update(info master.WorkerInfo) {
	r.Delete(info.GetId())
	r.Push(info)
}
