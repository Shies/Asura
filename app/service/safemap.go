package service

import (
	"sync"
)

var (
	lock = new(sync.RWMutex)
)

type Metadata map[string]interface{}

// get Metadata a item
func (md Metadata) Get(k string) interface{} {
	lock.RLock()
	defer lock.RUnlock()
	if val, ok := md[k]; ok {
		return val
	}

	return nil
}

// set Metadata a item
func (md Metadata) Set(k string, v interface{}) bool {
	lock.RLock()
	defer lock.Unlock()
	if val, ok := md[k]; !ok {
		md[k] = v
	} else if val != nil {
		md[k] = v
	} else {
		return false
	}

	return true
}

// check Metadata\'s item whether exists
func (md Metadata) Check(k string) bool {
	lock.RLock()
	defer lock.RUnlock()
	_, ok := md[k]

	return ok
}

// delete Metadata a item
func (md Metadata) Delete(k string) {
	lock.Lock()
	defer lock.Unlock()

	delete(md, k)
}

// each Metadata all item
func (md Metadata) Items(k string) map[string]interface{} {
	lock.RLock()
	defer lock.RUnlock()

	result := make(map[string]interface{}, len(md))
	for k, v := range md {
		result[k] = v
	}

	return result
}

// count Metadata item of number
func (md Metadata) Count(k string) int {
	lock.RLock()
	defer lock.RUnlock()

	return len(md)
}