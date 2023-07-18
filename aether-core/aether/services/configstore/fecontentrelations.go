// Services > Configstore > Content Relations

// This package is the content equivalent of user relations, it handles actions like subscribing and unsubscribing to a board, following a thread, and so on.

// This package also handles the optional SFWList for boards, which is used to highlight especially interesting communities. Neither this SFWList, nor the API for it is *not* a part of the protocol, only a part of this specific c0 client app.

/**
 *
 * Heads up - the way to use this is through configstore and you have to access this through GetContentRelations, and when done, you should do SetContentRelations, otherwise it won't be saved permanently.
 *
 * If you want to just read, you can read without doing a Get. But if you want to read and write, you should do get, edit and set, because that's the only way to retain changes.
 *
 */

package configstore

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Board struct {
	Fingerprint string
	Notify      bool
	LastSeen    int64
}

type Thread struct {
	Fingerprint string
	Notify      bool
	LastSeen    int64
}

type ContentRelations struct {
	lock          sync.Mutex
	Initialised   bool
	SubbedBoards  []Board
	SubbedThreads []Thread
	SFWList       sfwlist
}

func (c *ContentRelations) Init() {
	c.Initialised = true
}

// func (c *ContentRelations) GetAllSubbedThreads() []Thread {
// 	c.lock.Lock()
// 	defer c.lock.Unlock()
// 	return c.SubbedThreads
// }

/*----------  Subscription status  ----------*/

func (c *ContentRelations) IsSubbedBoard(fp string) (isSubbed, notifyEnabled bool, lastSeen int64) {
	loc := c.FindBoard(fp)
	if loc != -1 {
		return true, c.SubbedBoards[loc].Notify, c.SubbedBoards[loc].LastSeen
	}
	return false, false, 0
}

func (c *ContentRelations) FindBoard(fp string) int {
	for key, _ := range c.SubbedBoards {
		if c.SubbedBoards[key].Fingerprint == fp {
			return key
		}
	}
	return -1
}

func (c *ContentRelations) FindThread(fp string) int {
	for key, _ := range c.SubbedThreads {
		if c.SubbedThreads[key].Fingerprint == fp {
			return key
		}
	}
	return -1
}

func (c *ContentRelations) GetAllSubbedBoards() []Board {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.SubbedBoards
}

/*----------  Signals (silenced/notify, etc.) status  ----------*/

// SetBoardSignal sets the board signal into the storage. If a board is subscribed, we set the notify signal as well, if a subscription is removed, we remove the entry.
func (c *ContentRelations) SetBoardSignal(
	fp string, subscribed, notify bool, lastseen int64, lastSeenOnly bool) (committed bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if lastSeenOnly {
		c.insertLastSeenForBoard(fp, lastseen)
		return
	}
	if subscribed {
		c.insertBoard(fp, notify, lastseen, lastSeenOnly)
	} else {
		c.removeBoard(fp)
	}
	return true
}

/*----------  Internal work functions  ----------*/

func (c *ContentRelations) insertBoard(fp string, notify bool, lastseen int64, lastSeenOnly bool) {
	if i := c.FindBoard(fp); i != -1 {
		c.SubbedBoards[i].Notify = notify
		if lastseen > c.SubbedBoards[i].LastSeen {
			c.SubbedBoards[i].LastSeen = lastseen
		}
		return
	}
	c.SubbedBoards = append(c.SubbedBoards,
		Board{Fingerprint: fp, Notify: notify})
}

func (c *ContentRelations) insertLastSeenForBoard(fp string, lastseen int64) {
	if i := c.FindBoard(fp); i != -1 {
		if lastseen > c.SubbedBoards[i].LastSeen {
			c.SubbedBoards[i].LastSeen = lastseen
		}
	}
}

func (c *ContentRelations) insertThread(fp string, notify bool) {
	if i := c.FindThread(fp); i != -1 {
		c.SubbedThreads[i].Notify = notify
		return
	}
	c.SubbedThreads = append(c.SubbedThreads,
		Thread{Fingerprint: fp, Notify: notify})
}

func (c *ContentRelations) removeBoard(fp string) {
	if i := c.FindBoard(fp); i != -1 {
		c.SubbedBoards = append(c.SubbedBoards[0:i], c.SubbedBoards[i+1:len(c.SubbedBoards)]...)
	}
}

func (c *ContentRelations) removeThread(fp string) {
	if i := c.FindThread(fp); i != -1 {
		c.SubbedThreads = append(c.SubbedThreads[0:i], c.SubbedThreads[i+1:len(c.SubbedThreads)]...)
	}
}

// func (c *ContentRelations) SubBoard(fp string, notify bool) {
// 	c.insertBoard(fp, notify)
// }

// func (c *ContentRelations) UnsubBoard(fp string) {
// 	// c.lock.Lock()
// 	// defer c.lock.Unlock()
// 	c.removeBoard(fp)
// }

// func (c *ContentRelations) SubThread(fp string, notify bool) {
// 	c.lock.Lock()
// 	defer c.lock.Unlock()
// 	c.insertThread(fp, notify)
// }

// func (c *ContentRelations) UnsubThread(fp string) {
// 	c.lock.Lock()
// 	defer c.lock.Unlock()
// 	c.removeThread(fp)
// }

/*----------  SFW list  ----------*/

type sfwlist struct {
	lock       sync.Mutex
	LastUpdate int64
	Source     string
	Boards     []string
}

func (list *sfwlist) Update() {
	list.lock.Lock()
	defer list.lock.Unlock()
	if fc.GetSFWListDisabled() == false {
		return
	}
	if len(list.Source) == 0 {
		list.Source = "https://sfwlist.getaether.net/sfwlist.json"
	}
	// If last update was less than an hour ago, return
	// Check the sfwlist source
	list.commit()
}

func (list *sfwlist) GetSFWListDisabled() bool {
	return fc.GetSFWListDisabled()
}

func (list *sfwlist) SetSFWListDisabled(state bool) {
	fc.SetSFWListDisabled(state)
}

var LastSFWListUpdateInThisRun int64

func (list *sfwlist) Refresh() {
	list.lock.Lock()
	defer list.lock.Unlock()
	if list.GetSFWListDisabled() {
		return
	}
	if LastSFWListUpdateInThisRun != 0 {
		if time.Since(time.Unix(list.LastUpdate, 0)).Minutes() < 60 {
			// If it's been updated in the last 60 minutes, no point in refreshing.
			return
		}
	}
	list.LastUpdate = time.Now().Unix()
	// ^ We set the last update timestamp even if it fails - so that we'll only check every hour. If this wasn't the case, every action that would trigger a refresh after a failed call would trigger another call.
	if len(list.Source) == 0 {
		list.Source = "https://sfwlist.getaether.net/sfwlist.json"
	}
	list.Boards = []string{}
	response, err := http.Get(list.Source)
	if err != nil {
		log.Printf("Attempting to refresh the SFW list encountered an error. Err: %#v", err)
		return
	}
	defer response.Body.Close()
	resp, err2 := ioutil.ReadAll(response.Body)
	if err2 != nil {
		log.Printf("Attempting to read the SFW list after fetch encountered an error. Err: %#v", err2)
		return
	}
	rsl := RemoteSFWList{}
	err3 := json.Unmarshal(resp, &rsl)
	if err3 != nil {
		log.Printf("Attempting to parse the JSON of the SFW list encountered an error. Err: %#v", err2)
		return
	}
	log.Printf("SFW list arrived: %#v", rsl)
	if len(rsl.Boards) == 0 {
		return
	}
	list.Boards = rsl.Boards
	fc.SetContentRelations(fc.GetContentRelations())
	list.commit()
	LastSFWListUpdateInThisRun = time.Now().Unix()
}

func (list *sfwlist) IsSFWListedBoard(fp string) (isSFWListed bool) {
	list.Refresh()
	return list.FindBoardInSFWList(fp) != -1
}

func (list *sfwlist) FindBoardInSFWList(fp string) int {
	for key, _ := range list.Boards {
		if list.Boards[key] == fp {
			return key
		}
	}
	return -1
}

func (list *sfwlist) commit() {
	fc.SetContentRelations(fc.GetContentRelations())
}

/*----------  SFW List struct  ----------*/

type RemoteSFWList struct {
	Boards []string `json:"boards"`
}
