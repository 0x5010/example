package main

import (
	"bytes"
	"encoding/gob"
	"sync"

	"github.com/weaveworks/mesh"
)

type state struct {
	mutex sync.RWMutex
	set   map[mesh.PeerName]int
	self  mesh.PeerName
}

var _ mesh.GossipData = &state{}

func newState(self mesh.PeerName) *state {
	return &state{
		set:  map[mesh.PeerName]int{},
		self: self,
	}
}

func (s *state) get() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	c := 0
	for _, v := range s.set {
		c += v
	}
	return c
}

func (s *state) incr() *state {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.set[s.self]++
	return &state{
		set: s.set,
	}
}

func (s *state) copy() *state {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return &state{
		set: s.set,
	}
}

func (s *state) Encode() [][]byte {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(s.set); err != nil {
		panic(err)
	}
	return [][]byte{buf.Bytes()}
}

func (s *state) Merge(other mesh.GossipData) mesh.GossipData {
	return s.mergeComplete(other.(*state).copy().set)
}

func (s *state) mergeComplete(set map[mesh.PeerName]int) mesh.GossipData {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for peer, v := range set {
		if v > s.set[peer] {
			s.set[peer] = v
		}
	}
	return &state{
		set: s.set,
	}
}

func (s *state) mergeReceived(set map[mesh.PeerName]int) mesh.GossipData {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for peer, v := range set {
		if v <= s.set[peer] {
			delete(set, peer)
		} else {
			s.set[peer] = v
		}

	}
	return &state{
		set: set,
	}
}

func (s *state) mergeDelta(set map[mesh.PeerName]int) mesh.GossipData {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for peer, v := range set {
		if v <= s.set[peer] {
			delete(set, peer)
		} else {
			s.set[peer] = v
		}
	}

	if len(set) <= 0 {
		return nil
	}
	return &state{
		set: set,
	}
}
