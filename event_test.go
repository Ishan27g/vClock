package vClock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEventVector(t *testing.T) {
	leader := NewEventVector()

	// event from p1
	e1 := newEvent("event1", EventClock{
		"p1": 1,
		"p2": 1,
		"p3": 0,
	})
	// event 2 clock > event 1 clock
	e2 := newEvent("event2", EventClock{
		"p1": 1,
		"p2": 2,
		"p3": 0,
	})
	// event 3 is the latest message
	e3 := newEvent("event3", EventClock{
		"p1": 3,
		"p2": 3,
		"p3": 3,
	})
	// receive in wrong order
	leader.MergeEvent(Event{EventId: "event3", EventClock: e3.EventClock})
	leader.MergeEvent(Event{EventId: "event1", EventClock: e1.EventClock})
	leader.MergeEvent(Event{EventId: "event2", EventClock: e2.EventClock})

	order := leader.GetEventsOrder()
	assert.Equal(t, "event1", order[0])
	assert.Equal(t, "event2", order[1])
	assert.Equal(t, "event3", order[2])
}
