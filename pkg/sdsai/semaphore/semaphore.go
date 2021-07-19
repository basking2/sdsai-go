package semaphore

import (
	"sync"
	"sync/atomic"
)

type Semaphore struct {
	level int32
	cond  *sync.Cond
}

func NewSemaphore(level int32) *Semaphore {
	m := sync.Mutex{}
	c := sync.NewCond(&m)
	return &Semaphore{
		level: level,
		cond:  c,
	}
}

func (s *Semaphore) Down(diff int32) {
	s.cond.L.Lock()
	var v int32
	for v = atomic.LoadInt32(&s.level); v < diff; v = atomic.LoadInt32(&s.level) {
		// Wait for an up.
		s.cond.Wait()
	}

	atomic.StoreInt32(&s.level, v-diff)

	s.cond.L.Unlock()
}

// Return true if the down succeeds.
func (s *Semaphore) TryDown(diff int32) bool {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()
	if v := atomic.LoadInt32(&s.level); v >= diff {
		atomic.StoreInt32(&s.level, v-diff)
		return true
	} else {
		return false
	}

}

func (s *Semaphore) Up(diff int32) {
	s.cond.L.Lock()
	atomic.AddInt32(&s.level, diff)
	s.cond.Signal()
	s.cond.L.Unlock()
}

func (s *Semaphore) GetLevel() int32 {
	return atomic.LoadInt32(&s.level)
}
