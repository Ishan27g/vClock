package vClock

import "github.com/emirpasic/gods/lists/arraylist"

type Events interface {
	NewEvent(eventIdOrHash string, v2 VClock)
	// GetEventsOrder returns the events ordered according to vector clock
	// for the events
	GetEventsOrder() (eventIdsOrHashes []string)
}

// Event Every event is saved with the vector clock of when that event was received
type Event struct {
	events map[string]*VClock
	eventList *arraylist.List
}
func (e *Event) GetEventsOrder() (eventIdsOrHashes []string) {
	e.sort()
	var happenedBefore []string
	e.eventList.Each(func(_ int, value interface{}) {
		e := value.(event)
		happenedBefore = append(happenedBefore, e.eventIdOrHash)
	})
	return happenedBefore
}
func (e *Event) sort() {
	e.eventList.Sort(func(a, b interface{}) int {
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

// NewEvent merges the current vector clock for this event with the received vector clock
func (e *Event) NewEvent(eventIdOrHash string, v2 VClock) {
	if e.events[eventIdOrHash] == nil{
		e.events[eventIdOrHash] = &v2
		e.eventList.Add(event{
			eventIdOrHash: eventIdOrHash,
			clock:         v2,
		})
		return
	}
	eventClock := *e.events[eventIdOrHash]
	if eventClock == nil {

	}
	e.mergeClock(eventIdOrHash, v2)
}
func (e *Event) mergeClock(eventIdOrHash string, v2 VClock) {
	eventClock := *e.events[eventIdOrHash]
	for addr, clock := range v2 {
		if eventClock[addr] < clock {
			eventClock[addr] = clock
		}
	}
	for addr, clock := range eventClock {
		if v2[addr] > clock {
			eventClock[addr] = clock
		}
	}
	e.events[eventIdOrHash] = &eventClock

	newList := arraylist.New()
	for it := e.eventList.Iterator(); it.Next(); {
		e:= it.Value().(event)
		if e.eventIdOrHash == eventIdOrHash{
			newList.Add(event{
					eventIdOrHash: eventIdOrHash,
					clock:         eventClock,
				})
		}else {
			newList.Add(e)
		}
	}
}

func NewEventVectorClock() Events{
	e := Event{
		events: map[string]*VClock{},
		eventList: arraylist.New(),
	}
	return &e
}