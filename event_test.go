package vClock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEventVector(t *testing.T) {
	t.Parallel()

	leader := NewEventVector("p0", []string{"p2","p3","p4"})

	// event from p1
	e1 := newEvent("event1", EventClock{
		"p1" : 1,
		"p2" : 1,
		"p3" : 0,
	})
	// event 2 clock > event 1 clock
	e2 := newEvent("event2", EventClock{
		"p1" : 1,
		"p2" : 2,
		"p3" : 0,
	})
	// event 3 is the latest message
	e3 := newEvent("event3", EventClock{
		"p1" : 3,
		"p2" : 3,
		"p3" : 3,
	})

	// receive in wrong order
	leader.MergeEventClock("event3", e3.eventClock)
	leader.MergeEventClock("event1", e1.eventClock)
	leader.MergeEventClock("event2", e2.eventClock)

	order := leader.GetEventsOrder()
	assert.Equal(t, "event1", order[0])
	assert.Equal(t, "event2", order[1])
	assert.Equal(t, "event3", order[2])
}
