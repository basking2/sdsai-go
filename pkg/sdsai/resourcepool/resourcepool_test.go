package resourcepool

import (
	"sync"
	"testing"
)

type IntMaker struct{}

var wg sync.WaitGroup = sync.WaitGroup{}

func (*IntMaker) Create() (interface{}, error) {
	i := 1
	return &i, nil
}

func (*IntMaker) Check(interface{}) error {
	return nil
}

func (*IntMaker) Destroy(interface{}) {
	wg.Done()
}

func TestResourcePool(t *testing.T) {

	wg.Add(2)

	pool, err := NewResourcePool(
		&IntMaker{},
		2,
		2,
		2,
	)
	if err != nil {
		t.Error(err)
	}

	r1, err := pool.GetResource()
	if err != nil {
		t.Error(err)
	}

	r2, err := pool.GetResource()
	if err != nil {
		t.Error(err)
	}

	pool.ReturnResource(r1)
	pool.ReturnResource(r2)

	r1, err = pool.GetResource()
	if err != nil {
		t.Error(err)
	}

	r2, err = pool.GetResource()
	if err != nil {
		t.Error(err)
	}

	pool.ReturnResource(r1)
	pool.ReturnResource(r2)

	wg.Wait()

}
