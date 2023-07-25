// Backend > Dispatch > Exclusions
// This subsystem ensures that we are not repeatedly connecting to the same nodes over and over again within short periods of time.

package dispatch

import (
	"aether-core/aether/io/api"
	"aether-core/aether/services/globals"
	"fmt"
	"sync"
	"time"
)

type dispatcherExclusions struct {
	lock           sync.Mutex
	Exclusions     map[string]time.Time
	LastMaintained int64
}

var dpe dispatcherExclusions
var reverseDpe dispatcherExclusions

func (d *dispatcherExclusions) Add(a api.Address) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.prepare()
	d.maintain()
	d.Exclusions[d.canonicaliseAddr(a)] = time.Now()
}

func (d *dispatcherExclusions) IsExcluded(a api.Address) bool {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.prepare()
	d.maintain()
	liveExpiry := globals.BackendConfig.GetDispatchExclusionExpiryForLiveAddress()
	staticExpiry := globals.BackendConfig.GetDispatchExclusionExpiryForStaticAddress()
	ts := d.Exclusions[d.canonicaliseAddr(a)]
	if ts.Unix() > 0 {
		if a.Type == 255 || a.Type == 254 || a.Type == 253 { // Static
			if time.Since(ts) > staticExpiry {
				return false // we can connect
			}
			return true // too recent, can't connect
		}
		if a.Type == 2 || a.Type == 3 || a.Type == 4 { // Live
			if time.Since(ts) > liveExpiry {
				return false // we can connect
			}
			return true // too recent, can't connect
		}
		return true // not a type of node we recognise, can't connect
	}
	return false // not previously present in the excl list, we can connect.
}

/*----------  Maintenance / service methods  ----------*/

func (d *dispatcherExclusions) prepare() {
	if d.Exclusions == nil {
		d.Exclusions = make(map[string]time.Time)
	}
}

func (d *dispatcherExclusions) maintain() {
	now := time.Now()
	maintenanceCutoff := now.Add(-1 * time.Hour).Unix()
	exclusionsCutoff := now.Add(-12 * time.Hour).Unix()
	if d.LastMaintained > maintenanceCutoff {
		for k := range d.Exclusions {
			if d.Exclusions[k].Unix() > exclusionsCutoff {
				delete(d.Exclusions, k)
			}
		}
		d.LastMaintained = now.Unix()
	}
}

func (d *dispatcherExclusions) canonicaliseAddr(a api.Address) string {
	return fmt.Sprintf("%s %s %s", a.Location, a.Sublocation, a.Port)
}
