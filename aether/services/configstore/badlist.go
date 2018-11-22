// Services > Configstore > Badlist

// This package handles the checking of newly-arriving or existing content against known bad fingerprints.

package configstore

import (
	"aether-core/aether/services/toolbox"
	"encoding/json"
	"errors"
	"fmt"
	cdir "github.com/shibukawa/configdir"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type BadBoard struct {
	Fingerprint string
	Timestamp   int64
}
type BadThread struct {
	BoardFingerprint string
	Fingerprint      string
	Timestamp        int64
}
type BadPost struct {
	BoardFingerprint  string
	ThreadFingerprint string
	ParentFingerprint string
	Fingerprint       string
	Timestamp         int64
}
type BadVote struct {
	BoardFingerprint  string
	ThreadFingerprint string
	TargetFingerprint string
	Fingerprint       string
	Timestamp         int64
}
type BadKey struct {
	Fingerprint string
	Timestamp   int64
}
type BadTruststate struct {
	TargetFingerprint string
	Fingerprint       string
	Timestamp         int64
}
type BadAddress struct {
	Location           string
	Sublocation        string
	Port               uint16
	ClientVersionMajor uint8
	ClientVersionMinor uint16
	ClientVersionPatch uint16
	ClientName         string
}

type Badlist struct {
	lock        sync.Mutex
	LastUpdate  int64
	Source      string
	Boards      map[string]BadBoard
	Threads     map[string]BadThread
	Posts       map[string]BadPost
	Votes       map[string]BadVote
	Keys        map[string]BadKey
	Truststates map[string]BadTruststate
	Addresses   map[string]BadAddress
}

type badlistPayload struct {
	Boards      map[string]BadBoard
	Threads     map[string]BadThread
	Posts       map[string]BadPost
	Votes       map[string]BadVote
	Keys        map[string]BadKey
	Truststates map[string]BadTruststate
	Addresses   map[string]BadAddress
}

var BadlistInstance = Badlist{
	Boards:      make(map[string]BadBoard),
	Threads:     make(map[string]BadThread),
	Posts:       make(map[string]BadPost),
	Votes:       make(map[string]BadVote),
	Keys:        make(map[string]BadKey),
	Truststates: make(map[string]BadTruststate),
	Addresses:   make(map[string]BadAddress)}

var LastBadlistUpdateInThisRun int64

func (list *Badlist) Refresh() {
	list.lock.Lock()
	defer list.lock.Unlock()
	list.refresh()
}

func (list *Badlist) refresh() {
	list.fillFromDisk()
	if LastBadlistUpdateInThisRun != 0 {
		if time.Since(time.Unix(list.LastUpdate, 0)).Minutes() < 60 {
			// If it's been updated in the last 60 minutes, no point in refreshing.
			return
		}
	}
	list.LastUpdate = time.Now().Unix()
	// ^ We set the last update timestamp even if it fails - so that we'll only check every hour. If this wasn't the case, every action that would trigger a refresh after a failed call would trigger another call.
	if len(list.Source) == 0 {
		list.Source = "https://static.getaether.net/Badlist/Latest/badlist.json"
	}
	response, err := http.Get(list.Source)
	if err != nil {
		log.Printf("Attempting to refresh the Badlist encountered an error. Err: %#v", err)
		return
	}
	defer response.Body.Close()
	resp, err2 := ioutil.ReadAll(response.Body)
	if err2 != nil {
		log.Printf("Attempting to read the Badlist after fetch encountered an error. Err: %#v", err2)
		return
	}
	remoteBadlist := badlistPayload{}
	err3 := json.Unmarshal(resp, &remoteBadlist)
	if err3 != nil {
		log.Printf("Attempting to parse the JSON of the Badlist encountered an error. Err: %#v", err2)
		return
	}
	log.Printf("Badlist arrived: %#v", remoteBadlist)
	list.refreshFromNewData(&remoteBadlist)
	list.saveToDisk()
}

// parseBadlist refreshes the active badlist cache that is being used, and it saves the BL to the disk.
func (list *Badlist) refreshFromNewData(remoteBadlist *badlistPayload) {
	list.Boards = remoteBadlist.Boards
	list.Threads = remoteBadlist.Threads
	list.Posts = remoteBadlist.Posts
	list.Votes = remoteBadlist.Votes
	list.Keys = remoteBadlist.Keys
	list.Truststates = remoteBadlist.Truststates
	list.Addresses = remoteBadlist.Addresses
	// if len(remoteBadlist.Boards) > 0 {
	// 	list.Boards = remoteBadlist.Boards
	// }
	// if len(remoteBadlist.Threads) > 0 {
	// 	list.Threads = remoteBadlist.Threads
	// }
	// if len(remoteBadlist.Posts) > 0 {
	// 	list.Posts = remoteBadlist.Posts
	// }
	// if len(remoteBadlist.Votes) > 0 {
	// 	list.Votes = remoteBadlist.Votes
	// }
	// if len(remoteBadlist.Keys) > 0 {
	// 	list.Keys = remoteBadlist.Keys
	// }
	// if len(remoteBadlist.Truststates) > 0 {
	// 	list.Truststates = remoteBadlist.Truststates
	// }
	// if len(remoteBadlist.Addresses) > 0 {
	// 	list.Addresses = remoteBadlist.Addresses
	// }
	LastBadlistUpdateInThisRun = time.Now().Unix()
}

/*----------  Read from disk  ----------*/

func (list *Badlist) fillFromDisk() error {
	binI, orgI, appI := getIdentifiers()
	configDirs := cdir.New(orgI, appI)
	folder := configDirs.QueryFolderContainsFile(filepath.Join(binI, "badlist.json"))
	if folder != nil {
		listJson, _ := folder.ReadFile(filepath.Join(binI, "badlist.json"))
		err := json.Unmarshal(listJson, list)
		if err != nil || fmt.Sprintf("%#v", string(listJson)) == "\"{}\"" {
			return errors.New(fmt.Sprintf("Badlist file is corrupted. Please fix the file, or delete it. If deleted a new badlist will be generated with default values. Error: %#v, ListJson: %#v", err, string(listJson)))
		}
	}
	// Folder is nil - this file does not exist. We fill nothing. We write an empty file as a starting point.
	list.saveToDisk()
	return nil
}

/*----------  Save to disk  ----------*/

func (list *Badlist) SaveToDisk() error {
	list.lock.Lock()
	defer list.lock.Unlock()
	err := list.saveToDisk()
	return err
}

/*
saveToDisk saves the file to memory. This is usually called after a Set operation. This works for both backend and frontend, which means in a normal desktop system, there will be two copies of this file. A small inefficiency to retain the independence of both binaries.
*/
func (list *Badlist) saveToDisk() error {
	listAsByte, err := json.Marshal(list)
	if err != nil {
		log.Fatal(fmt.Sprintf("JSON marshaler encountered an error while marshaling this Badlist into JSON. List: %#v, Error: %#v", list, err))
	}
	binI, orgI, appI := getIdentifiers()
	configDirs := cdir.New(orgI, appI)
	folders := configDirs.QueryFolders(cdir.Global)
	toolbox.CreatePath(folders[0].Path)
	toolbox.CreatePath(filepath.Join(folders[0].Path, binI))
	writeAheadPath := filepath.Join(folders[0].Path, binI, "badlist_writeahead.json")
	targetPath := filepath.Join(folders[0].Path, binI, "badlist.json")
	err2 := ioutil.WriteFile(writeAheadPath, listAsByte, 0644)
	if err2 != nil {
		return err2
	}
	err3 := os.Rename(writeAheadPath, targetPath)
	if err3 != nil {
		return err3
	}
	return nil
}

/*----------  Util functions  ----------*/

func getIdentifiers() (string, string, string) {
	binI := "backend"
	orgI, appI := Btc.OrgIdentifier, Btc.AppIdentifier
	if len(appI) == 0 || len(orgI) == 0 {
		orgI, appI = Ftc.OrgIdentifier, Ftc.AppIdentifier
		binI = "frontend"
	}
	return binI, orgI, appI
}

/*----------  Query functions  ----------*/

/*

	This cascades.
	We can't afford to do a full graph traversal here, because that will get exponentially costly.

	What we do is this:

-----------------
	For a board to be bad,
	- its fp in the bad list
	- its owners fp in the bad list
-----------------
	For a thread to be bad, it has
	- its fp in the bad list
	- its owners fp in the bad list
	- its boards fp in the bad list
	X - its boards owner fp in the bad list

	^ We don't check for this last one - because we'd have to query its board, and find its owner, and then do an owner query check. This gets very expensive very fast.

	Mind that these checks are done in the ingestion phase where a node is syncing and saving data from another node. This has to be fast, it cannot do cross-queries.

	This means we cannot query the owner of a parent in a content tree. That means if the parent is bad because it is in the badlist, that goes out, but if it's bad because its key is bad, it won't go out, and it will be accepted in. However, if a parent is bad due to its key, then the parent itself won't be able to come in, and all content that considers it an ancestor will be invisible, thus still serving our purpose. It just adds an inefficiency in rejecting content, it does not make that content visible.

	In the case this is abused to fill databases, we can always make this a cross-checking scanner, and start to fully discard the children as well, not a big deal.
-----------------
	For a post to be bad, it has
	- its fp in the bad list
	- its owners fp in the bad list

	X - its parent post's fp in the bad list
	X - its parent post's owner fp in the bad list

	- its threads fp in the bad list
	X - its threads owner fp in the bad list

	- its boards fp in the bad list
	X - its boards owner fp in the bad list
-----------------
	For a vote to be bad, it has
	- its fp in the bad list
	- its owners fp in the bad list

	X - its target fp in the bad list
	X - its target owner fp in the bad list

	- its threads fp in the bad list
	X - its threads owner fp in the bad list

	- its boards fp in the bad list
	X - its boards owner fp in the bad list
-----------------
	For a key to be bad, it has
	- its fp in the bad list
-----------------
	For a truststate to be bad, it has
	- its fp in the bad list
	- its owners fp in the bad list

	X - its target fp in the bad list
	X - its target owner fp in the bad list
-----------------
	For an address to be bad, it has
	- its loc/subloc/port in the bad list

*/

func (list *Badlist) IsBadBoard(fp, ownerfp string) bool {
	list.lock.Lock()
	defer list.lock.Unlock()
	if LastBadlistUpdateInThisRun == 0 {
		list.refresh()
	}
	return list.isBadBoard(fp, ownerfp)
}

func (list *Badlist) isBadBoard(fp, ownerfp string) bool {
	// Short exits
	if len(list.Boards) == 0 && len(list.Keys) == 0 {
		return false
	}
	// Full check path
	if o := list.Boards[fp]; o.Fingerprint == fp {
		return true
	}
	if o := list.Keys[ownerfp]; o.Fingerprint == ownerfp {
		return true
	}
	return false
}

func (list *Badlist) IsBadThread(fp, boardfp, ownerfp string) bool {
	list.lock.Lock()
	defer list.lock.Unlock()
	if LastBadlistUpdateInThisRun == 0 {
		list.refresh()
	}
	return list.isBadThread(fp, boardfp, ownerfp)
}

func (list *Badlist) isBadThread(fp, boardfp, ownerfp string) bool {
	// Short exits
	if len(list.Threads) == 0 && len(list.Boards) == 0 && len(list.Keys) == 0 {
		return false
	}
	// Full check path
	if o := list.Threads[fp]; o.Fingerprint == fp {
		return true
	}
	if o := list.Keys[ownerfp]; o.Fingerprint == ownerfp {
		return true
	}
	return list.isBadBoard(boardfp, ownerfp)
}

func (list *Badlist) IsBadPost(fp, boardfp, threadfp, parentfp, ownerfp string) bool {
	list.lock.Lock()
	defer list.lock.Unlock()
	if LastBadlistUpdateInThisRun == 0 {
		list.refresh()
	}
	return list.isBadPost(fp, boardfp, threadfp, parentfp, ownerfp)
}

func (list *Badlist) isBadPost(fp, boardfp, threadfp, parentfp, ownerfp string) bool {
	// Short exits
	if len(list.Posts) == 0 && len(list.Threads) == 0 && len(list.Boards) == 0 && len(list.Keys) == 0 {
		return false
	}
	// Full check path
	if o := list.Posts[fp]; o.Fingerprint == fp {
		return true
	}
	if o := list.Keys[ownerfp]; o.Fingerprint == ownerfp {
		return true
	}
	/*
		Special to post, we are not checking whether the *parent post* of this post is valid or not - we're just checking whether the post itself is, its thread is, and its post is. Why? Because checking for the parent creates a recursive cascade, and that can get very expensive very fast.

		As of now, if a post is not present, the posts under it will be orphaned and not visible, even if they're in the database - and that achieves the goal.
	*/
	return list.isBadThread(threadfp, boardfp, ownerfp)
}

func (list *Badlist) IsBadVote(fp, boardfp, threadfp, targetfp, ownerfp string) bool {
	list.lock.Lock()
	defer list.lock.Unlock()
	if LastBadlistUpdateInThisRun == 0 {
		list.refresh()
	}
	return list.isBadVote(fp, boardfp, threadfp, targetfp, ownerfp)
}

// This is not even called as of now.
func (list *Badlist) isBadVote(fp, boardfp, threadfp, targetfp, ownerfp string) bool {
	// Short exits
	if len(list.Votes) == 0 && len(list.Posts) == 0 && len(list.Threads) == 0 && len(list.Boards) == 0 && len(list.Keys) == 0 {
		return false
	}
	// Full check path
	if o := list.Votes[fp]; o.Fingerprint == fp {
		return true
	}
	if o := list.Keys[ownerfp]; o.Fingerprint == ownerfp {
		return true
	}
	/*
		Special to vote, we're just checking whether the parent thread and parent board is valid. If we go into checking the posts and that cascades, it gets very expensive very fast.
	*/
	return list.isBadThread(threadfp, boardfp, ownerfp)
}

func (list *Badlist) IsBadKey(fp string) bool {
	list.lock.Lock()
	defer list.lock.Unlock()
	if LastBadlistUpdateInThisRun == 0 {
		list.refresh()
	}
	return list.isBadKey(fp)
}

func (list *Badlist) isBadKey(fp string) bool {
	// Short exits
	if len(list.Keys) == 0 {
		return false
	}
	// Full check path
	if o := list.Keys[fp]; o.Fingerprint == fp {
		return true
	}
	return false
}

func (list *Badlist) IsBadTruststate(fp, targetfp, ownerfp string) bool {
	list.lock.Lock()
	defer list.lock.Unlock()
	if LastBadlistUpdateInThisRun == 0 {
		list.refresh()
	}
	return list.isBadTruststate(fp, targetfp, ownerfp)
}

func (list *Badlist) isBadTruststate(fp, targetfp, ownerfp string) bool {
	// Short exits
	if len(list.Truststates) == 0 && len(list.Keys) == 0 {
		return false
	}
	// Full check path
	if o := list.Truststates[fp]; o.Fingerprint == fp {
		return true
	}
	if o := list.Keys[ownerfp]; o.Fingerprint == ownerfp {
		return true
	}
	return list.isBadKey(targetfp)
}

func (list *Badlist) IsBadAddress(loc, subloc string, port uint16) bool {
	list.lock.Lock()
	defer list.lock.Unlock()
	if LastBadlistUpdateInThisRun == 0 {
		list.refresh()
	}
	return list.isBadAddress(loc, subloc, port)
}

func (list *Badlist) isBadAddress(loc, subloc string, port uint16) bool {
	// Short exits
	if len(list.Addresses) == 0 {
		return false
	}
	// Full check path
	parsed := fmt.Sprintf("%s:%s/%s", loc, port, subloc)
	addr := list.Addresses[parsed]
	if addr.Location == loc && addr.Port == port && addr.Sublocation == subloc {
		return true
	}
	return false
}
