// Frontend > FEStructs > Notifications
// This library provides the buckets for the notifications subsystem based on the signals coming from the frontend refresher. This will be saved into the KvInstance. It also maintains that KvInstance state, and provides handles for opening and closing it.

package festructs

import (
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/toolbox"
	"fmt"
	"sort"
	"sync"
	"time"
)

/*
  The operating principle here is that we have a notifications container for each of the entities that are self created. This is nice, because these containers are automatically created.

  When we are compiling the posts, we get the delta, and we stick that delta into the notifications system. This system gets the self posts, creates the buckets for it, and of the stuff that ends up being actually responses, puts them into the appropriate buckets.
*/

type NotificationsCarrier struct {
	lock         sync.RWMutex
	instantiated bool
	Id           int `storm:"id"` // always 1, this is a singleton.
	Containers   map[string]NotificationsContainer
	LastSeen     int64
}

var NotificationsSingleton NotificationsCarrier

type NotificationsContainer struct {
	Thread               CompiledThread
	Post                 CompiledPost
	LastUpdate           int64
	Muted                bool
	NotificationsBuckets []NotificationsBucket
}

type NotificationsBucket struct {
	LastUpdate         int64
	ResponsePosts      map[string]int64         // [Fingerprint]Timestamp
	ResponsePostsUsers map[string]CUserUsername // [Fingerprint]Username
	Read               bool
}

func NewNotificationsBucket(now int64) NotificationsBucket {
	return NotificationsBucket{
		LastUpdate:         now,
		ResponsePosts:      make(map[string]int64),
		ResponsePostsUsers: make(map[string]CUserUsername),
	}
}

func (c *NotificationsContainer) Insert(ce CompiledPost, now int64) {
	// Check if post exists anywhere. We might have raised a notification for it already.
	for k, _ := range c.NotificationsBuckets {
		if c.NotificationsBuckets[k].ResponsePosts[ce.Fingerprint] != 0 {
			return
		}
	}
	// Go through all notifications buckets available that aren't already read.
	var latestNonReadNBLastUpdate int64
	latestNonReadNBIndex := -1
	for k, _ := range c.NotificationsBuckets {
		if c.NotificationsBuckets[k].Read {
			// If bucket is read, pass, we can't do anything to that any more.
			continue
		}
		if c.NotificationsBuckets[k].LastUpdate > latestNonReadNBLastUpdate {
			// At the end of the loop, we'll have the most recent not-yet-read notification bucket. There should always be one, generally speaking.
			latestNonReadNBLastUpdate = c.NotificationsBuckets[k].LastUpdate
			latestNonReadNBIndex = k
		}
	}
	postUsernameStruct := ce.Owner.GetUsername()
	if latestNonReadNBIndex != -1 {
		// We have a not-yet-read notification bucket we can insert into
		c.NotificationsBuckets[latestNonReadNBIndex].LastUpdate = now
		c.NotificationsBuckets[latestNonReadNBIndex].ResponsePosts[ce.Fingerprint] = now
		c.NotificationsBuckets[latestNonReadNBIndex].ResponsePostsUsers[ce.Fingerprint] = postUsernameStruct
		return
	}
	// We have no notification bucket to house this. Create a new one.
	nb := NewNotificationsBucket(now)
	nb.ResponsePosts[ce.Fingerprint] = now
	nb.ResponsePostsUsers[ce.Fingerprint] = postUsernameStruct
	c.NotificationsBuckets = append(c.NotificationsBuckets, nb)
}

/*----------  Instantiation / uninstantiation  ----------*/

// ResetNotificationsSingleton uninstantiates and reinstantiates the notifications singleton.
func ReinstantiateNotificationsSingleton() {
	NotificationsSingleton = NewNotificationsCarrier()
	InstantiateNotificationsSingleton()
}

func InstantiateNotificationsSingleton() {
	if NotificationsSingleton.instantiated {
		return
	}
	logging.Logf(3, "Single read happens in InstantiateNotificationsSingleton>One")
	err := globals.KvInstance.One("Id", 1, &NotificationsSingleton)
	if err != nil {
		logging.Logf(1, "Fetching data from KvStore of instantiating the NotificationsCarrier had an error. Error: %v", err)
		// Attempt to recreate.
		nc := NewNotificationsCarrier()
		NotificationsSingleton = nc
	}
	NotificationsSingleton.instantiated = true
}

func NewNotificationsCarrier() NotificationsCarrier {
	logging.Logf(1, "Reinitialising notifications container as blank.")
	ncNew := NotificationsCarrier{
		Id:         1,
		Containers: make(map[string]NotificationsContainer),
		// instantiated: true,
	}
	return ncNew
}

func (nc *NotificationsCarrier) save() {
	logging.Logf(3, "Save happens in NotificationsCarrier>Save")
	globals.KvInstance.Save(nc)
}

func (nc *NotificationsCarrier) Save() {
	nc.lock.Lock()
	defer nc.lock.Unlock()
	nc.save()
}

func (nc *NotificationsCarrier) SaveAndUninstantiate() {
	nc.lock.Lock()
	defer nc.lock.Unlock()
	nc.instantiated = false
	nc.save()
}

func (nc *NotificationsCarrier) MarkSeen() {
	nc.lock.Lock()
	defer nc.lock.Unlock()
	nc.LastSeen = time.Now().Unix()
}

func (nc *NotificationsCarrier) MarkRead(fp string) {
	nc.lock.Lock()
	defer nc.lock.Unlock()
	container := nc.Containers[fp]
	for k, _ := range container.NotificationsBuckets {
		container.NotificationsBuckets[k].Read = true
	}
	nc.Containers[fp] = container
}

func (nc *NotificationsCarrier) markAllAsRead() {
	for i, _ := range nc.Containers {
		for j, _ := range nc.Containers[i].NotificationsBuckets {
			nc.Containers[i].NotificationsBuckets[j].Read = true
		}
	}
}

func (nc *NotificationsCarrier) MarkAllAsRead() {
	nc.lock.Lock()
	defer nc.lock.Unlock()
	nc.markAllAsRead()
}

/*----------  Maintenance  ----------*/

func (nc *NotificationsCarrier) Prune() {
	cutoff := toolbox.CnvToCutoffDays(30)
	for k, _ := range nc.Containers {
		if nc.Containers[k].LastUpdate < cutoff {
			delete(nc.Containers, k)
		}
	}
	logging.Logf(1, "Notifications prune complete.")
}

/*----------  Insertion and mark read/unread  ----------*/

func (nc *NotificationsCarrier) InsertPosts(posts []CompiledPost) {
	nc.lock.Lock()
	defer nc.lock.Unlock()
	var nonSelfPosts []CompiledPost
	now := time.Now().Unix()
	// If it's a self post, add to the map
	for k, _ := range posts {
		if posts[k].SelfCreated {
			nContainer := nc.Containers[posts[k].Fingerprint]
			nContainer.Post = posts[k]
			nContainer.LastUpdate = now
			nc.Containers[posts[k].Fingerprint] = nContainer
			continue
		}
		nonSelfPosts = append(nonSelfPosts, posts[k])
	}
	// ^ Be mindful that we're removing self posts from the lists to be checked. That means responding to yourself will not raise a notification. Neat.
	// If not a self post, check if its parent matches a known self thread or post.
	for k, _ := range nonSelfPosts {
		if nc.responseToSelfPost(&nonSelfPosts[k]) {
			logging.Logf(2, "This is a response to a self post! %v", nonSelfPosts[k].Fingerprint)
			// It's a response to a self post. Insert it.
			nContainer := nc.Containers[nonSelfPosts[k].Parent]
			nContainer.Insert(nonSelfPosts[k], now)
			nc.Containers[nonSelfPosts[k].Parent] = nContainer
			continue
		}
		if nc.responseToSelfThread(&nonSelfPosts[k]) {
			// It's a response to a self thread.
			nContainer := nc.Containers[nonSelfPosts[k].Thread]
			nContainer.Insert(nonSelfPosts[k], now)
			nc.Containers[nonSelfPosts[k].Thread] = nContainer
			continue
		}
		// ^ Be mindful of the order. We are inserting the notification into the closest parent - if this is a response to a self post that was response to a self thread, it will be shown as a notification that says it's a response to the self post.
	}
}

// Heads up, for this to actually be useful, the insert thread needs to happen before insert posts, so that the posts will be able to check for existence of this self thread.
func (nc *NotificationsCarrier) InsertThreads(threads []CompiledThread) {
	nc.lock.Lock()
	defer nc.lock.Unlock()
	for k, _ := range threads {
		if !threads[k].SelfCreated {
			continue
		}
		// If it's a self thread, add to the map
		nContainer := nc.Containers[threads[k].Fingerprint]
		nContainer.Thread = threads[k]
		nContainer.LastUpdate = time.Now().Unix()
		// logging.Logf(1, "This is the notifications container that's crashing as we try to insert: %#v", nc)
		nc.Containers[threads[k].Fingerprint] = nContainer
		// ^ We don't have a use for a thread that is not a self here, so those are completely ignored.
	}
}

func (nc *NotificationsCarrier) responseToSelfPost(ce *CompiledPost) bool {
	if len(nc.Containers[ce.Parent].Post.Fingerprint) == 0 {
		// The usual case. Nonexistent, meaning this isn't one of the user's self posts.
		return false
	}
	if nc.Containers[ce.Parent].Post.Fingerprint != ce.Parent {
		// The weird case. The only way this would happen is by programming mistake: the map key and the map value's fingerprint do not match. This is impossible unless there's a bug, but this is a defence for it nevertheless.
		return false
	}
	return true
}

func (nc *NotificationsCarrier) responseToSelfThread(ce *CompiledPost) bool {
	if len(nc.Containers[ce.Thread].Thread.Fingerprint) == 0 {
		return false
	}
	if nc.Containers[ce.Thread].Thread.Fingerprint != ce.Thread {
		return false
	}
	return true
}

const (
	REPLY_TO_THREAD = 1
	REPLY_TO_POST   = 2
)

/*----------  Listification to send to client  ----------*/
type CompiledNotification struct {
	Type                    int // REPLY_TO_THREAD, REPLY_TO_POST
	Text                    string
	ResponsePosts           []string
	ResponsePostsUsers      map[string]CUserUsername
	ParentThread            CompiledThread
	ParentPost              CompiledPost
	CreationTimestamp       int64
	NewestResponseTimestamp int64
	Read                    bool
}

func (c *CompiledNotification) GenerateText() {
	switch c.Type {
	case REPLY_TO_POST:
		c.generateReplyToPostText()
	case REPLY_TO_THREAD:
		c.generateReplyToThreadText()
		// c.Text = generateReplyToThrea3dText(len(c.ResponsePosts), c.ParentThread)
	default:
		logging.Logf(1, "Compiled notification has an unknown type. CompiledNotification: %v", c)
	}
}

func (c *CompiledNotification) generateReplyToPostText() {
	var parentPost = c.ParentPost
	var respCount = len(c.ResponsePosts)
	var shortenedPostBody string
	if len([]rune(parentPost.Body)) < 48 {
		shortenedPostBody = parentPost.Body
	} else {
		shortenedPostBody = fmt.Sprintf("%s…", string([]rune(parentPost.Body)[0:48]))
	}
	if respCount == 1 {
		// un := c.ResponsePostsUsers[c.ResponsePosts[0]]
		c.Text = fmt.Sprintf("replied to post “%s”", shortenedPostBody)
		return
	}
	c.Text = fmt.Sprintf("%d replies to post “%s”", respCount, shortenedPostBody)
}

func (c *CompiledNotification) generateReplyToThreadText() {
	var parentThread = c.ParentThread
	var respCount = len(c.ResponsePosts)
	var shortenedThrName string
	if len([]rune(parentThread.Name)) < 48 {
		shortenedThrName = parentThread.Name
	} else {
		shortenedThrName = fmt.Sprintf("%s…", string([]rune(parentThread.Name)[0:48]))
	}
	if respCount == 1 {
		c.Text = fmt.Sprintf("replied to thread “%s”", shortenedThrName)
		return
	}
	c.Text = fmt.Sprintf("%d replies to thread “%s”", respCount, shortenedThrName)
}

type CNotificationsList []CompiledNotification

// Listify is the logic that runs every time there is a need to send the client the notifications that we have now.
func (nc *NotificationsCarrier) Listify() (CNotificationsList, int64) {
	nc.lock.Lock()
	defer nc.lock.Unlock()
	start := time.Now()

	/*====================================================================
	=            Notifications silence gate on search reindex            =
	====================================================================*/
	if globals.FrontendTransientConfig.SilenceNotificationsOnce == true {
		logging.Logf(2, "SilenceNotificationsOnce is present, we're silencing this listify and the next one won't be silent.")
		if len(nc.Containers) > 0 {
			/*
				This is gated at containers > 0 because we do not want to gate our notifications in the case this is a zero-out entity generation.

				The interesting thing there is that you might legitimately have zero notifications, and if you receive a notification while this is present, i.e. in the same run, that notification would get silenced.

				This would only happen if ALL of these are true:

				a) You deleted your search index, but not your frontend config.

				b) You have legitimately *never* received any notification in the app before (because if you did, that will be muted, and the mute gate will be removed).

				c) You receive your first ever notification right after this happens and this gate is triggered on a new, actually unseen notification, silencing it without you seeing it.

				d) The app is not restarted before the first notification comes. (Restart also extinguishes this gate on its own.)

				I think it's pretty far fetched, honestly, simply because deleting the index is not going to be a common occurrence by itself.
			*/
			globals.FrontendTransientConfig.SilenceNotificationsOnce = false
			nc.markAllAsRead()
			nc.save()
		}
	}
	/*=====  End of Notifications silence gate on search reindex  ======*/

	cnl := CNotificationsList{}
	// For each container
	for k, _ := range nc.Containers {
		// Skip if muted
		if nc.Containers[k].Muted {
			continue
		}
		var nType int
		nType = REPLY_TO_POST // Post by default, if thread, flip it
		if thr := len(nc.Containers[k].Thread.Fingerprint); thr > 0 {
			nType = REPLY_TO_THREAD
		}
		// For every bucket in container
		for k2, _ := range nc.Containers[k].NotificationsBuckets {
			if len(nc.Containers[k].NotificationsBuckets[k2].ResponsePosts) == 0 {
				// Skip if no responses (generally impossible, but good to guard against)
				continue
			}
			// Convert fingerprints:lastupdate map to []string fp
			var rpFps []string

			rpUsers := nc.Containers[k].NotificationsBuckets[k2].ResponsePostsUsers
			for k, _ := range nc.Containers[k].NotificationsBuckets[k2].ResponsePosts {
				rpFps = append(rpFps, k)
			}
			// Figure out the newest response in the bucket and use its timestamp
			var newest int64
			for k3, _ := range nc.Containers[k].NotificationsBuckets[k2].ResponsePosts {
				if nc.Containers[k].NotificationsBuckets[k2].ResponsePosts[k3] > newest {
					newest = nc.Containers[k].NotificationsBuckets[k2].ResponsePosts[k3]
				}
			}
			// Create the compiled notification object
			cn := CompiledNotification{
				Type:                    nType,
				ResponsePosts:           rpFps,
				ResponsePostsUsers:      rpUsers,
				ParentPost:              nc.Containers[k].Post,
				ParentThread:            nc.Containers[k].Thread,
				CreationTimestamp:       nc.Containers[k].NotificationsBuckets[k2].LastUpdate,
				Read:                    nc.Containers[k].NotificationsBuckets[k2].Read,
				NewestResponseTimestamp: newest,
			}
			cn.GenerateText()
			// Add to our main bucket and return
			cnl = append(cnl, cn)
		}
		// Within the given box, we surface the ones that contain the newest responses higher.
		sort.Slice(cnl, func(i, j int) bool {
			// If less or more, fairly straightforward. Sort by that.
			if cnl[i].NewestResponseTimestamp > cnl[j].NewestResponseTimestamp {
				return true
			}
			if cnl[i].NewestResponseTimestamp < cnl[j].NewestResponseTimestamp {
				return false
			}
			// If the values are the same:
			// 1: We look at the parent post or thread creation (not last update, because we want the first time the notification was raised, so as to not let people 'bump' notifications in other people's timelines by updating the response continuously)
			var parentCreationI, parentCreationJ int64
			if cnl[i].ParentPost.Creation > cnl[i].ParentThread.Creation {
				parentCreationI = cnl[i].ParentPost.Creation
			} else {
				parentCreationI = cnl[i].ParentThread.Creation
			}
			if cnl[j].ParentPost.Creation > cnl[j].ParentThread.Creation {
				parentCreationJ = cnl[j].ParentPost.Creation
			} else {
				parentCreationJ = cnl[j].ParentThread.Creation
			}
			if parentCreationI > parentCreationJ {
				return true
			}
			if parentCreationI < parentCreationJ {
				return false
			}
			// If still the same ...
			// 2: We look at the body text and thread text. So long as those aren't the same, the sort order will be deterministic.
			a := fmt.Sprintf("%s%s%s", cnl[i].ParentPost.Body, cnl[i].ParentThread.Name, cnl[i].ParentThread.Body)
			b := fmt.Sprintf("%s%s%s", cnl[j].ParentPost.Body, cnl[j].ParentThread.Name, cnl[j].ParentThread.Body)
			return a > b
		})
	}
	// logging.Logf(1, "Notifications buckets: %#v", nc.Containers)
	// logging.Logf(1, "Notifications being sent after compile: %#v", cnl)
	elapsed := time.Since(start)
	logging.Logf(1, "Notifications listification took %v", elapsed)
	return cnl, nc.LastSeen
}
