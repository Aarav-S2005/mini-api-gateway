package store

import "sync/atomic"

type Registry struct {
	snap atomic.Pointer[Snapshot]
}

func NewRegistry(s *Snapshot) *Registry {
	r := &Registry{}
	r.snap.Store(s)
	return r
}

func (r *Registry) Get() *Snapshot { return r.snap.Load() }

func (r *Registry) Store(s *Snapshot) { r.snap.Store(s) }
