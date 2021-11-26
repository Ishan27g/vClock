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
	Get() VClock

	// SendEvent returns the current vector clock after first updating self clock and then
	// updating the individual clocks for corresponding addresses
	SendEvent(eventIdOrHash string, addresses []string) VClock

	// ReceiveEvent updates the current vector clock using element wise maximum with the passed vector clock
	ReceiveEvent(eventIdOrHash string, v VClock)

	// GetEventsOrder returns the events ordered according to vector clock
	// for the events
	GetEventsOrder() (eventIdsOrHashes []string)

	Clear()
	print()
}

// VClock is a map of peer-address and its individual clock
type VClock map[string]int

// Every event is saved with the vector clock of when that event was received
type event struct {
	eventIdOrHash string
	clock VClock
}

type vClock struct {
	lock sync.Mutex
	self string
	vectorClock VClock
	addressList *arraylist.List
	events *arraylist.List
}

func (v *vClock) Clear() {
	v.lock.Lock()
	v.vectorClock = nil
	v.addressList.Clear()
	v.addressList = nil
	v.events = nil
	v.lock.Unlock()
}

func (v *vClock) print() {
	fmt.Println(v.events.String())
	// for a, c := range v.vectorClock {
	//	fmt.Println(a, c)
	// }
}

// Get returns the current vector clock
func (v *vClock) Get() VClock {
	return v.vectorClock
}


// GetEventsOrder returns the events ordered according to vector clock
// for the events
func (v *vClock) GetEventsOrder() (eventIdsOrHashes []string) {
	v.sortEvents()
	var happenedBefore []string
	v.events.Each(func(_ int, value interface{}) {
		e := value.(event)
		happenedBefore = append(happenedBefore, e.eventIdOrHash)
	})
	return happenedBefore
}

// sortEvents sorts the events according to happened before relation
func (v *vClock)sortEvents(){
	v.events.Sort(func(a, b interface{}) int {
		e1 := a.(event)
		e2 := a.(event)
		c1 := compareClock(e1.clock, e2.clock)
		c2 := compareClock(e2.clock, e1.clock)
		if c1 && c2 { // both are same
			return 0
		} else if c1 && !c2 { // e1 happened before
			return -1
		} else {
			return 1
		}
	})
}

// SendEvent returns the current vector clock after updating the individual clocks for these entries
func (v *vClock) SendEvent(eventIdOrHash string, address []string) VClock{
	v.lock.Lock()
	defer v.lock.Unlock()
	// update the individual clock entry for self
	v.event(v.self)

	for _, a := range address {
		v.event(a)
	}
	v.events.Add(event{
		eventIdOrHash: eventIdOrHash,
		clock:         v.Get(),
	})
	return v.Get()
}
// ReceiveEvent updates the current vector clock using element wise maximum with the passed vector clock
func (v *vClock) ReceiveEvent(eventIdOrHash string, v1 VClock) {
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
	v.events.Add(event{
		eventIdOrHash: eventIdOrHash,
		clock:         v.Get(),
	})
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

func Init(self string, peers []string) VectorClock {
	v := vClock{
		lock:        sync.Mutex{},
		vectorClock: make(map[string]int),
		self: self,
		addressList: arraylist.New(),
		events: arraylist.New(),
	}

	for _, peer := range peers {
		v.initClock(peer)
	}
	v.initClock(v.self)
	v.addressList.Sort(utils.StringComparator)
	return &v
}

func (v *vClock)initClock(peer string) {
	v.vectorClock[peer] = 0
	v.addressList.Add(peer)
}

// compareClock returns true if v1 is before or concurrent to v2
func compareClock(v1 VClock, v2 VClock) bool {
	if len(v1) != len(v2) {
		return false
	}
	v1IsBefore := true
	for addr, v1Clock := range v1 {
		v2Clock := v2[addr]
		if v2Clock == 0 && v1Clock != 0{
			continue
		}
		if v2Clock < v1Clock {
			v1IsBefore = false
		}
	}
	for addr, v2Clock := range v2 {
		v1Clock := v1[addr]
		if v1Clock == 0 && v2Clock != 0{
			continue
		}
		if v2Clock < v1Clock {
			v1IsBefore = false
		}
	}
	return v1IsBefore
}
