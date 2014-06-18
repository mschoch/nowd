package main

import (
	"sync"
	"time"
)

type TimeRevCache struct {
	mutex      sync.RWMutex
	expiration time.Duration
	revs       map[string]int
	times      map[string]time.Time
	values     map[string]interface{}
}

func NewTimeRevCache(expiration time.Duration) *TimeRevCache {
	return &TimeRevCache{
		expiration: expiration,
		revs:       make(map[string]int),
		times:      make(map[string]time.Time),
		values:     make(map[string]interface{}),
	}
}

func (t TimeRevCache) CheckAndUpdate(key string, rev int, value interface{}) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	currentRev, existsAndNotExpired := t.getrev(key)
	if existsAndNotExpired {
		if rev <= currentRev {
			return false
		}
		t.set(key, rev, value)
		return true
	}
	t.set(key, rev, value)
	return true
}

func (t TimeRevCache) set(key string, rev int, value interface{}) {
	delete(t.revs, key)
	delete(t.times, key)
	delete(t.values, key)
	t.revs[key] = rev
	t.times[key] = time.Now()
	t.values[key] = value
}

func (t TimeRevCache) getrev(key string) (int, bool) {
	currentRev, exists := t.revs[key]
	if exists {
		// check expiration
		createdTime, exists := t.times[key]
		if exists {
			lifetime := time.Since(createdTime)
			if lifetime <= t.expiration {
				return currentRev, true
			}
			delete(t.revs, key)
			delete(t.times, key)
			return 0, false
		}
		delete(t.revs, key)
		return 0, false
	}
	return 0, false
}

func (t TimeRevCache) Expire() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for key, createdTime := range t.times {
		lifetime := time.Since(createdTime)
		if lifetime > t.expiration {
			delete(t.revs, key)
			delete(t.times, key)
			delete(t.values, key)
		}
	}
}

func (t TimeRevCache) Values() map[string]interface{} {
	rv := make(map[string]interface{})
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	for k, v := range t.values {
		rv[k] = v
	}

	return rv
}
