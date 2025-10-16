package userbalance

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type userLock struct {
	mu       *sync.Mutex
	lastUsed time.Time
}

type lockGroup struct {
	mu    sync.Mutex
	locks map[uuid.UUID]*userLock
}

var (
	lockTTL         = 5 * time.Minute
	cleanupInterval = 1 * time.Minute
)

func (lg *lockGroup) cleanup() {
	lg.mu.Lock()
	defer lg.mu.Unlock()

	now := time.Now()
	for id, l := range lg.locks {
		if now.Sub(l.lastUsed) > lockTTL {
			delete(lg.locks, id)
		}
	}
}

func (lg *lockGroup) GetOrCreate(id uuid.UUID) *sync.Mutex {
	lg.mu.Lock()
	defer lg.mu.Unlock()

	if l, exists := lg.locks[id]; exists {
		l.lastUsed = time.Now()
		return l.mu
	}

	newLock := &userLock{mu: &sync.Mutex{}, lastUsed: time.Now()}
	lg.locks[id] = newLock
	return newLock.mu
}
