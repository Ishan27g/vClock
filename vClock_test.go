package vClock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {

	m1 := make(EventClock)
	m1["p1"] = 1
	m1["p2"] = 1
	m1["p3"] = 1

	m2 := make(EventClock)
	m2["p1"] = 1
	m2["p2"] = 1
	m2["p3"] = 1

	m3 := make(EventClock)
	m3["p1"] = 1
	m3["p2"] = 2
	m3["p3"] = 1

	m4 := make(EventClock)
	m4["p1"] = 2
	m4["p2"] = 2
	m4["p3"] = 3

	assert.True(t, compareClock(m1, m2))
	assert.True(t, compareClock(m2, m3))
	assert.False(t, compareClock(m3, m2))
	assert.True(t, compareClock(m1, m4))
	assert.True(t, compareClock(m2, m4))
	assert.True(t, compareClock(m3, m4))
	assert.False(t, compareClock(m4, m3))

	m5 := make(EventClock)
	m5["p1"] = 20
	m5["p2"] = 20
	m5["p3"] = 40
	m5["p4"] = 20

	m6 := make(EventClock)
	m6["p1"] = 20
	m6["p2"] = 21

	assert.True(t, compareClock(m5, m6))
	assert.False(t, compareClock(m6, m5))
}

func TestSendEventEvent(t *testing.T) {
	p1 := Init("p1")
	p1.SendEvent("event", []string{"p2"})
	assert.Equal(t, 1, p1.Get()["p2"])
	p1.Reset()
}
func TestClocksMatch(t *testing.T) {

	p1 := Init("p1")
	p2 := Init("p2")

	// p1 sends vClock to p2
	p1.SendEvent("event", []string{"p2"})
	// p2 receives vClock from p1
	p2.ReceiveEvent("event", p1.Get())

	// both clocks match
	assert.Equal(t, p1.Get(), p2.Get())
	p1.Reset()
	p2.Reset()
}
