package semaphore

import (
	"sync"
	"testing"
)

func TestSemaphore(t *testing.T) {
	s := NewSemaphore(10)

	for i := 0; i < 10; i++ {
		s.Down(1)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	b := false

	go func() {
		s.Down(5)
		b = true
		s.Up(10)
		wg.Done()
	}()

	if b {
		t.Error("Failed.")
	}

	s.Up(5)
	s.Down(10)
	if !b {
		t.Error("Failed 2.")
	}

}

func TestSemaphoreTryDown(t *testing.T) {
	s := NewSemaphore(10)

	b := s.TryDown(10)
	if !b {
		t.Error("TryDown should succeed.")
	}

	b = s.TryDown(1)
	if b {
		t.Error("TryDown should fail.")
	}
}
