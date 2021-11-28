package vClock

import (
	"fmt"
	"sync"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/utils"
)
// VectorClock provides interface to implement vector clock
type VectorClock interface {
	// Get returns the current vector clock
	Get() EventClock

	// SendEvent returns the current vector clock after first updating self clock and then
	// updating the individual clocks for corresponding addresses
	SendEvent(eventIdOrHash string, addresses []string) EventClock

	// ReceiveEvent updates the current vector clock using element wise maximum with the passed vector clock
	ReceiveEvent(eventIdOrHash string, v EventClock)

	Clear()
	print()
}

// EventClock is a map of peer-address and its individual clock

type vClock struct {
	lock sync.Mutex
	self        string
	vectorClock EventClock
	addressList *arraylist.List
}

func (v *vClock) AddEvent(eventIdOrHash string, v1 EventClock) {
	v.lock.Lock()
	defer v.lock.Unlock()
	// merge with received clock
	for address, newClock := range v1 {
		if v.addressList.Contains(address){
			v.updateClock(address, newClock)
		}else { // if new address
			v.initClock(address)
			v.updateClock(address, newClock)
		}
	}

}

func (v *vClock) Clear() {
	v.lock.Lock()
	v.vectorClock = nil
	v.addressList.Clear()
	v.addressList = nil
	v.lock.Unlock()
}

func (v *vClock) print() {
	for a, c := range v.vectorClock {
		fmt.Println(a, c)
	}
}

// Get returns the current vector clock
func (v *vClock) Get() EventClock {
	return v.vectorClock
}

// SendEvent returns the current vector clock after updating the individual clocks for these entries
func (v *vClock) SendEvent(eventIdOrHash string, address []string) EventClock {
	v.lock.Lock()
	defer v.lock.Unlock()
	// update the individual clock entry for self
	v.event(v.self)

	for _, a := range address {
		v.event(a)
	}
	return v.Get()
}
// ReceiveEvent updates the current vector clock using element wise maximum with the passed vector clock
func (v *vClock) ReceiveEvent(eventIdOrHash string, v1 EventClock) {
	v.lock.Lock()
	defer v.lock.Unlock()
	// update local clock
	v.event(v.self)
	// merge with received clock
	for address, newClock := range v1 {
		if v.addressList.Contains(address){
			v.updateClock(address, newClock)
		}else { // if new address
			v.initClock(address)
			v.updateClock(address, newClock)
		}
	}
}
// event updates the individual clock entry for this entry
func (v *vClock) event(address string) {
	currentClock := v.vectorClock[address]
	v.vectorClock[address] = currentClock + 1
}
// updateClock updates the individual clock if it is lower than the new clock
func (v *vClock) updateClock(address string, newClock int) {
	if v.vectorClock[address] < newClock {
		v.vectorClock[address] = newClock
	}
}

func Init(self string) VectorClock {
	v := vClock{
		lock:        sync.Mutex{},
		vectorClock: make(map[string]int),
		self: self,
		addressList: arraylist.New(),
	}
	v.initClock(v.self)
	v.addressList.Sort(utils.StringComparator)
	return &v
}

func (v *vClock)initClock(peer string) {
	v.vectorClock[peer] = 0
	v.addressList.Add(peer)
}
