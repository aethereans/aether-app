// Services > Configstore > FEUserRelations

// This library handles the in/outs of the actions that the frontend user has taken. This is not used in the backend.

package configstore

import (
	"sync"
)

type User struct {
	Fingerprint string
	Domain      string
}
type UserRelations struct {
	lock            sync.Mutex
	Initialised     bool
	Following       BatchUser
	Blocked         BatchUser
	ModElected      BatchUser
	ModDisqualified BatchUser
}

func (u *UserRelations) Init() {
	u.Initialised = true
}

func (u *UserRelations) FollowUser(fp, domain string) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.Following.Insert(fp, domain)
}

func (u *UserRelations) UnfollowUser(fp, domain string) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.Following.Remove(fp, domain)
}

func (u *UserRelations) BlockUser(fp, domain string) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.Blocked.Insert(fp, domain)
}

func (u *UserRelations) UnblockUser(fp, domain string) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.Blocked.Remove(fp, domain)
}

func (u *UserRelations) ModElectUser(fp, domain string) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.ModElected.Insert(fp, domain)
}

func (u *UserRelations) UnModElectUser(fp, domain string) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.ModElected.Remove(fp, domain)
}

func (u *UserRelations) ModDisqualifyUser(fp, domain string) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.ModDisqualified.Insert(fp, domain)
}

// UnModDisqualify means reverting a mod-disqualify. it's a little confusing but consistent.
func (u *UserRelations) UnModDisqualifyUser(fp, domain string) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.ModDisqualified.Remove(fp, domain)
}

type BatchUser []User

func (b *BatchUser) Find(fp, domain string) int {
	for k := range *b {
		if (*b)[k].Fingerprint == fp && (*b)[k].Domain == domain {
			return k
		}
	}
	return -1
}

func (b *BatchUser) Insert(fp, domain string) {
	if i := b.Find(fp, domain); i != -1 {
		return
	}
	(*b) = append(*b, User{Fingerprint: fp, Domain: domain})
}

func (b *BatchUser) Remove(fp, domain string) {
	if i := b.Find(fp, domain); i != -1 {
		*b = append((*b)[0:i], (*b)[i+1:len(*b)]...)
	}
}
