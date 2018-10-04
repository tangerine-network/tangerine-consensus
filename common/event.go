// Copyright 2018 The dexon-consensus-core Authors
// This file is part of the dexon-consensus-core library.
//
// The dexon-consensus-core library is free software: you can redistribute it
// and/or modify it under the terms of the GNU Lesser General Public License as
// published by the Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.
//
// The dexon-consensus-core library is distributed in the hope that it will be
// useful, but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Lesser
// General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the dexon-consensus-core library. If not, see
// <http://www.gnu.org/licenses/>.

package common

import (
	"container/heap"
	"time"
)

type timeEventFn func(time.Time)

type timeEvent struct {
	t  time.Time
	fn timeEventFn
}

// timeEvents implements a Min-Heap structure.
type timeEvents []timeEvent

func (h timeEvents) Len() int           { return len(h) }
func (h timeEvents) Less(i, j int) bool { return h[i].t.Before(h[j].t) }
func (h timeEvents) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *timeEvents) Push(x interface{}) {
	*h = append(*h, x.(timeEvent))
}
func (h *timeEvents) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Event implements the Observer pattern.
type Event struct {
	timeEvents timeEvents
}

// NewEvent creates a new event instance.
func NewEvent() *Event {
	te := timeEvents{}
	heap.Init(&te)
	return &Event{
		timeEvents: te,
	}
}

// RegisterTime to get notified on and after specific time.
func (e *Event) RegisterTime(t time.Time, fn timeEventFn) {
	heap.Push(&e.timeEvents, timeEvent{
		t:  t,
		fn: fn,
	})
}

// NotifyTime and trigger function callback.
func (e *Event) NotifyTime(t time.Time) {
	if len(e.timeEvents) == 0 {
		return
	}
	for !t.Before(e.timeEvents[0].t) {
		te := heap.Pop(&e.timeEvents).(timeEvent)
		te.fn(t)
		if len(e.timeEvents) == 0 {
			return
		}
	}
}
