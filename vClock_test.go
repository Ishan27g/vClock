package vClock

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestSendEventEvent(t *testing.T) {
	p1 := Init("p1", []string{"p2","p3","p4"})
	p1.SendEvent("event", []string{"p2"})
	assert.Equal(t, 1, p1.Get()["p2"])
	p1.Clear()
}
func TestClocksMatch(t *testing.T) {

	p1 := Init("p1", []string{"p2","p3","p4"})
	p2 := Init("p2", []string{"p3","p1","p4"})

	// p1 sends vClock to p2
	p1.SendEvent("event", []string{"p2"})
	// p2 receives vClock from p1
	p2.ReceiveEvent("event", p1.Get())

	// both clocks match
	assert.Equal(t, p1.Get(), p2.Get())
	p1.Clear()
	p2.Clear()
}

func TestClocksMatchMultiple(t *testing.T) {
	t.Cleanup(func() {
		p1 := Init("p1", []string{"p2","p3"})
		p2 := Init("p2", []string{"p3","p1"})
		p3 := Init("p3", []string{"p1","p4"})

		p1p2 := make(chan EventClock)
		p1p3 := make(chan EventClock)
		p2p1 := make(chan EventClock)
		p2p3 := make(chan EventClock)

		wg := sync.WaitGroup{}
		wg.Add(8)

		// p1 waits till it receives vClock from p2 when p2 sends it
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			p1.ReceiveEvent("event2", <- p2p1)
		}(&wg)
		// p2 waits till it receives vClock from p1 when p1 sends it
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			p2.ReceiveEvent("event1", <- p1p2)
		}(&wg)

		// process 3 waits till it receives vClock from p1
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			p3.ReceiveEvent("event2", <- p1p3)
		}(&wg)

		// process 3 waits till it receives vClock from p2
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			p3.ReceiveEvent("event1", <- p2p3)
		}(&wg)

		// process 1
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			// p1 receives user-event & sends vClock to p2 & p3
			p1.SendEvent("event1", []string{"p2", "p3"})
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				p1p2 <- p1.Get()
			}(wg)
			p1p3 <- p1.Get()

		}(&wg)

		// process 2
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			// p2 receives user-event & sends vClock to p3 & p1
			p2.SendEvent("event2", []string{"p3", "p1"})
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				p2p3 <- p2.Get()
			}(wg)
			p2p1 <- p2.Get()

		}(&wg)

		wg.Wait()
		 //<- time.After(1 * time.Second)
		// all clocks match
		assert.Equal(t, p1.Get(), p2.Get())
		// common entries match
		assert.Equal(t, p1.Get()["p1"], p3.Get()["p1"])
		assert.Equal(t, p1.Get()["p3"], p3.Get()["p3"])

		p1.Clear()
		p2.Clear()
		p3.Clear()
	})


}
func TestCompare(t *testing.T){

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
}
