package vClock

import (
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
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
	leader.MergeEvent(cloudEvent("event3", e3.EventClock))
	leader.MergeEvent(cloudEvent("event1", e1.EventClock))
	leader.MergeEvent(cloudEvent("event2", e2.EventClock))

	order := leader.GetEventsOrder()
	assert.Equal(t, "event1", order[0])
	assert.Equal(t, "event2", order[1])
	assert.Equal(t, "event3", order[2])
}
func cloudEvent(id string, data EventClock) cloudevents.Event {
	return convertToCloud(event{
		EventId:    id,
		EventClock: data,
	})
}
