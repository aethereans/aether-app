// Services > ConfigStore > Neighbours
// This module maintains a list of nodes that we want to keep syncing with.

/*
What does this do?

This is where we keep the nodes we want to keep syncing with. Since our initial sync is relatively expensive, we want to still be regularly checking in with nodes that we've synced before.

How does this work?

This is a fixed size (configurable) stack that always pops out the oldest, and pushes in the newest. Effectively, every time you complete a sync, you push in a new node into this. Every time you need a node that you've synced before, you pop one out from this. If you do a pop, a spacer will be inserted at the beginning to keep it 10 entities still.

Example 1:

10 9 8 7 6 5 4 3 2 1 -- > 1 pops out

S 10 9 8 7 6 5 4 3 2 is the result.

When you pop out something and the connection fails, you don't reinsert, which means the spacer will remain and move forward.

When there is something that pushes without popping from it first, that means it is going to be inserted via removing the oldest spacer, or if no spacers are present, the oldest entity.

Example 2:

10 9 8 7 S S S 3 2 1

Insert 11

11 10 9 8 7 S S 3 2 1
              ^ spacer here got eaten

Example 3:

If there are no spacers, the oldest entry is gone

11 10 9 8 7 6 5 4 3 2
                    ^ 1 here got eaten

Why?

Assume that you're making a sync call to your neighbours every 1 minute. And you're connecting to a new node every 10 minutes. This structure will allow you to both keep connecting to nodes you know (by popping and pushing them into the stack) and keep the list updated via not pushing back nodes that go offline. It also allows the processes that do not pop from this stack to introduce a new neighbour without creating complex communication paths between processes and functions.

*/

package configstore

import (
	"sync"
)

// Structs
type Address struct {
	Location    string
	Sublocation string
	Port        uint16
}

type NeighboursList struct {
	lock       sync.Mutex
	Neighbours []Address
}

// Internal helpers

func (m *NeighboursList) indexOf(a Address) int {
	// Why reverse? Because if there are multiple, we want the oldest, which will be closer to the right side. This is useful because we want to delete older spacer sometimes.
	for i := len(m.Neighbours) - 1; i >= 0; i-- {
		if m.Neighbours[i].Location == a.Location &&
			m.Neighbours[i].Sublocation == a.Sublocation &&
			m.Neighbours[i].Port == a.Port {
			return i
		}
	}
	return -1
}

func (m *NeighboursList) removeByIndex(idx int) {
	first := m.Neighbours[:idx]
	second := m.Neighbours[idx+1 : len(m.Neighbours)]
	m.Neighbours = append(first, second...)
}

func (m *NeighboursList) insert(a Address) {
	m.Neighbours = append([]Address{a}, m.Neighbours...)
}

func (m *NeighboursList) removeOldest() {
	if len(m.Neighbours) > 0 {
		m.removeByIndex(len(m.Neighbours) - 1)
	}
}

func (m *NeighboursList) removeOldestSpacer() {
	oldestSpacerIndex := m.indexOf(Address{})
	if oldestSpacerIndex != -1 { // there is a spacer we can remove
		m.removeByIndex(oldestSpacerIndex)
	}
}

func isSpacer(a Address) bool {
	return a.Location == "" && a.Sublocation == "" && a.Port == 0
}

func (m *NeighboursList) insertSpacer() {
	m.insert(Address{})
}

func (m *NeighboursList) fillIfNeeded() {
	currentStackLen := len(m.Neighbours)
	if currentStackLen < bc.GetNeighbourCount() {
		for i := currentStackLen; i < bc.GetNeighbourCount(); i++ {
			m.insert(Address{})
		}
	}
}

func (m *NeighboursList) pruneIfNeeded() {
	for len(m.Neighbours) > bc.GetNeighbourCount() {
		m.removeOldest()
	}
}

// Main API

func (m *NeighboursList) Push(loc, subloc string, port uint16) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.fillIfNeeded()
	// if there are less than set stack size, fill with spacers.
	m.removeOldestSpacer()
	// if there are any spacers, remove the oldest, so we can fill this one in.
	m.insert(Address{Location: loc, Sublocation: subloc, Port: port})
	// insert into stack
	m.pruneIfNeeded()
	// bring back the size to bc.GetNeighbourCount() if it goes overboard.
}

func (m *NeighboursList) Pop() (string, string, uint16) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.fillIfNeeded()
	// Send the last item, and remove it from the list.
	ejectedItem := m.Neighbours[len(m.Neighbours)-1]
	m.Neighbours = m.Neighbours[:len(m.Neighbours)-1]
	// insert spacer at the beginning.
	m.insertSpacer()
	return ejectedItem.Location, ejectedItem.Sublocation, ejectedItem.Port
}
