/**
 * @Author: Resynz
 * @Date: 2022/1/6 10:42
 */
package config

import "sync"

type TaskMap struct {
	*sync.RWMutex
	Map map[int64]bool
}

func (s *TaskMap) Exists(id int64) bool {
	s.RLock()
	defer s.RUnlock()
	_,ok:=s.Map[id]
	if !ok {
		return false
	}
	return true
}

func (s *TaskMap) Set(id int64)  {
	s.Lock()
	defer s.Unlock()
	s.Map[id] = true
}

func (s *TaskMap) Delete(id int64)  {
	s.Lock()
	defer s.Unlock()
	delete(s.Map,id)
}
