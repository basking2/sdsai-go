package resourcepool

import (
	"time"
)

// How to move a resource through its lifetime.
// Create, Check, and Destroy.
type ResourceManager interface {
	// Create a resource and return it.
	Create() (interface{}, error)

	// Check if a resource should be destroyed.
	Check(interface{}) error

	// Destroy a resource.
	Destroy(interface{}) error
}

// A resource record.
type Resource struct {
	// The number of times this resource was fetched from the pool.
	Uses int

	// The time.Now().Unix() for when this object was created.
	CreatedAt int64

	// The created resource.
	Resource interface{}

	// The pool this resource came from.
	Pool *ResourcePool
}

type ResourcePool struct {
	MaxUses         int
	MaxAge          int64
	ResourceManager ResourceManager
	FreeResources   chan *Resource
	CreateResources chan int
}

func (r *Resource) Close() error {
	return r.Pool.ReturnResource(r)
}

func (r *Resource) Destroy() error {
	return r.Pool.DestroyResource(r)
}

func (pool *ResourcePool) UseResource(f func(*Resource) error) error {
	if r, e := pool.GetResource(); e == nil {
		defer pool.ReturnResource(r)
		return f(r)
	} else {
		return e
	}
}

func (pool *ResourcePool) GetResource() (*Resource, error) {
	select {
	case v := <-pool.FreeResources:
		v.Uses += 1
		return v, nil

	case <-pool.CreateResources:
		if r, e := pool.ResourceManager.Create(); e == nil {
			return &Resource{
				Uses:      1,
				CreatedAt: time.Now().Unix(),
				Pool:      pool,
				Resource:  r,
			}, nil
		} else {
			return nil, e
		}
	}
}

func (pool *ResourcePool) DestroyResource(r *Resource) error {
	e := pool.ResourceManager.Destroy(r.Resource)

	pool.CreateResources <- 1

	return e
}

func (pool *ResourcePool) ReturnResource(r *Resource) error {
	if pool.MaxUses > 0 && r.Uses >= pool.MaxUses {
		return pool.DestroyResource(r)
	} else if pool.MaxAge > 0 && time.Now().Unix()-r.CreatedAt > pool.MaxAge {
		return pool.DestroyResource(r)
	} else {
		pool.FreeResources <- r
		return nil
	}
}

func NewResourcePool(
	resourceManager ResourceManager,
	maxInstances int,
	maxUses int,
	maxAge int64,
) (*ResourcePool, error) {

	if maxInstances <= 0 {
		panic("Resource pools require a positive limit of instances.")
	}

	pool := ResourcePool{
		ResourceManager: resourceManager,
		MaxUses:         maxUses,
		MaxAge:          maxAge,
		FreeResources:   make(chan *Resource, maxInstances),
		CreateResources: make(chan int, maxInstances),
	}

	for i := 0; i < maxInstances; i++ {
		pool.CreateResources <- 1
	}

	return &pool, nil
}
