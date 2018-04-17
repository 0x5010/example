package main

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/weaveworks/mesh"
)

type peer struct {
	st      *state
	send    mesh.Gossip
	actions chan<- func()
	quit    chan struct{}
	logger  *log.Logger
}

var _ mesh.Gossiper = &peer{}

func newPeer(self mesh.PeerName, logger *log.Logger) *peer {
	actions := make(chan func())
	p := &peer{
		st:      newState(self),
		send:    nil,
		actions: actions,
		quit:    make(chan struct{}),
		logger:  logger,
	}
	go p.loop(actions)
	return p
}

func (p *peer) loop(actions <-chan func()) {
	for {
		select {
		case f := <-actions:
			f()
		case <-p.quit:
			return
		}
	}
}

func (p *peer) register(send mesh.Gossip) {
	p.actions <- func() { p.send = send }
}

func (p *peer) get() int {
	return p.st.get()
}

func (p *peer) incr() int {
	res := 0
	c := make(chan struct{})
	p.actions <- func() {
		defer close(c)
		st := p.st.incr()
		if p.send != nil {
			p.send.GossipBroadcast(st)
		} else {
			p.logger.Printf("no sender configured; not broadcasting update right now")
		}
		res = st.get()
	}
	<-c
	return res
}

func (p *peer) stop() {
	close(p.quit)
}

// Gossip 返回完整状态的副本
func (p *peer) Gossip() mesh.GossipData {
	complete := p.st.copy()
	p.logger.Printf("Gossip => complete %v", complete.set)
	return complete
}

// OnGossip 合并 gossip data(buf)到state中
// 返回修改的部分
func (p *peer) OnGossip(buf []byte) (mesh.GossipData, error) {
	var set map[mesh.PeerName]int
	if err := gob.NewDecoder(bytes.NewReader(buf)).Decode(&set); err != nil {
		return nil, err
	}

	delta := p.st.mergeDelta(set)
	if delta == nil {
		p.logger.Printf("OnGossip %v => delta %v", set, delta)
	} else {
		p.logger.Printf("OnGossip %v => delta %v", set, delta.(*state).set)
	}
	return delta, nil
}

// OnGossipBroadcast 合并 gossip data(buf)到state中
// 返回修改的部分
func (p *peer) OnGossipBroadcast(src mesh.PeerName, buf []byte) (mesh.GossipData, error) {
	var set map[mesh.PeerName]int
	if err := gob.NewDecoder(bytes.NewReader(buf)).Decode(&set); err != nil {
		return nil, err
	}

	received := p.st.mergeReceived(set)
	if received == nil {
		p.logger.Printf("OnGossipBroadcast %s %v => delta %v", src, set, received)
	} else {
		p.logger.Printf("OnGossipBroadcast %s %v => delta %v", src, set, received.(*state).set)
	}
	return received, nil
}

// OnGossipUnicast 合并 gossip data(buf)到state中
func (p *peer) OnGossipUnicast(src mesh.PeerName, buf []byte) error {
	var set map[mesh.PeerName]int
	if err := gob.NewDecoder(bytes.NewReader(buf)).Decode(&set); err != nil {
		return err
	}

	complete := p.st.mergeComplete(set)
	p.logger.Printf("OnGossipUnicast %s %v => complete %v", src, set, complete)
	return nil
}
