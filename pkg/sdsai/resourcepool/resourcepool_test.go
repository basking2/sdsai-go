package resourcepool

import (
	"sync"
	"testing"
)

type IntMaker struct {
	Wg sync.WaitGroup
}

func (*IntMaker) Create() (interface{}, error) {
	i := 1
	return &i, nil
}

func (*IntMaker) Check(interface{}) error {
	return nil
}

func (i *IntMaker) Destroy(interface{}) error {
	i.Wg.Done()
	return nil
}

func TestResourcePool(t *testing.T) {

	im := IntMaker{}

	im.Wg.Add(2)

	pool, err := NewResourcePool(
		&im,
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

	r1.Close()
	r2.Close()

	r1, err = pool.GetResource()
	if err != nil {
		t.Error(err)
	}

	r2, err = pool.GetResource()
	if err != nil {
		t.Error(err)
	}

	r1.Close()
	r2.Close()

	im.Wg.Wait()
}

func TestResourcePoolLoad(t *testing.T) {

	itrs := 1000

	im := IntMaker{}

	im.Wg.Add(itrs)

	pool, err := NewResourcePool(
		&im,
		itrs,
		1,
		-1,
	)
	if err != nil {
		t.Error(err)
	}

	println(pool)

	for i := 0; i < itrs; i++ {
		go func() {
			r, _ := pool.GetResource()
			go func() {
				r.Close()
			}()
		}()
	}

	im.Wg.Wait()

}
