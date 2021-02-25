package degrade

import (
	store2 "Asura/app/middleware/degrade/store"
	blade "Asura/src"
)

// Cache is the abstract struct for any cache impl
type Cache struct {
	store store2.Store
}

// Filter is used to check is cache required for every request
type Filter func(*blade.Context) bool

// Policy is used to abstract different cache policy
type Policy interface {
	Key(*blade.Context) string
	Handler(store2.Store) blade.HandlerFunc
}

// New will create a new Cache struct
func New(store store2.Store) *Cache {
	c := &Cache{
		store: store,
	}
	return c
}

// Cache is used to mark path as customized cache policy
func (c *Cache) Cache(policy Policy, filter Filter) blade.HandlerFunc {
	return func(ctx *blade.Context) {
		if filter != nil && !filter(ctx) {
			return
		}
		policy.Handler(c.store)(ctx)
	}
}
