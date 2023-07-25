// Frontend > KeyValueStore > Find
// This file takes some search IDs and it converts them to actual results we can deliver to the user.

package kvstore

import (
	"aether-core/aether/frontend/beapiconsumer"
	"aether-core/aether/frontend/festructs"
	"aether-core/aether/frontend/search"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
)

/*
  Search: Get the IDs where the search term matches.
  Find: Find the actual entities the IDs point to.

*/

/*----------  Search & find posts  ----------*/

type Searchable interface {
	GetEntityType() string
	GetBoardFp() string
	GetThreadFp() string
	GetParentFp() string
	GetFingerprint() string
	GetUserFp() string
}

// Searchable instead of the concrete type because this is shared with the new feed
func findPost(sId Searchable) (festructs.CompiledPost, error) {
	c := festructs.ThreadCarrier{}
	logging.Logf(3, "Single read happens in findPost>One")
	err := globals.KvInstance.One("Fingerprint", sId.GetThreadFp(), &c)
	if err != nil {
		logging.Logf(1, "We could not find the thread of this post. SearchId: %v", sId)
		return festructs.CompiledPost{}, err
	}
	i := c.Posts.Find(sId.GetFingerprint())
	if i == -1 {
		return festructs.CompiledPost{}, nil
	}
	return c.Posts[i], nil
}

func findPosts(searchResults search.SearchResults) []festructs.CompiledPost {
	var posts []festructs.CompiledPost

	for k := range searchResults.Results {
		if searchResults.Results[k].Id.EntityType != "Post" {
			continue
		}
		p, err := findPost(&searchResults.Results[k].Id)
		if err != nil {
			continue
		}
		posts = append(posts, p)
	}
	return posts
}

func SearchPosts(searchText string) (festructs.CPostBatch, search.ScoreMap, error) {
	resp, err := search.Search(searchText)
	if err != nil {
		return []festructs.CompiledPost{}, make(search.ScoreMap), err
	}
	return festructs.CPostBatch(findPosts(resp)), makeScoreMap(resp), nil
}

/*----------  Search & find threads  ----------*/
// Searchable instead of the concrete type because this is shared with the new feed
func findThread(sId Searchable) (festructs.CompiledThread, error) {
	c := festructs.ThreadCarrier{}
	logging.Logf(3, "Single read happens in findThread>One")
	err := globals.KvInstance.One("Fingerprint", sId.GetFingerprint(), &c)
	if err != nil {
		logging.Logf(1, "We could not find this thread. SearchId: %v", sId)
		return festructs.CompiledThread{}, err
	}
	return c.Threads[0], nil
}

func findThreads(searchResults search.SearchResults) []festructs.CompiledThread {
	var threads []festructs.CompiledThread

	for k := range searchResults.Results {
		if searchResults.Results[k].Id.EntityType != "Thread" {
			continue
		}
		p, err := findThread(&searchResults.Results[k].Id)
		if err != nil {
			continue
		}
		threads = append(threads, p)
	}
	return threads
}

func SearchThreads(searchText string) (festructs.CThreadBatch, search.ScoreMap, error) {
	resp, err := search.Search(searchText)
	if err != nil {
		return []festructs.CompiledThread{}, make(search.ScoreMap), err
	}
	return festructs.CThreadBatch(findThreads(resp)), makeScoreMap(resp), nil
}

/*----------  Search & find boards  ----------*/

func findBoard(searchResult search.SearchResult) (festructs.CompiledBoard, error) {
	searchId := searchResult.Id
	c := festructs.BoardCarrier{}
	logging.Logf(3, "Single read happens in findBoard>One")
	err := globals.KvInstance.One("Fingerprint", searchId.Fingerprint, &c)
	if err != nil {
		logging.Logf(1, "We could not find this board. SearchId: %v", searchId)
		return festructs.CompiledBoard{}, err
	}
	return c.Boards[0], nil
}

func findBoards(searchResults search.SearchResults) []festructs.CompiledBoard {
	var boards []festructs.CompiledBoard

	for k := range searchResults.Results {
		if searchResults.Results[k].Id.EntityType != "Board" {
			continue
		}
		p, err := findBoard(searchResults.Results[k])
		if err != nil {
			continue
		}
		boards = append(boards, p)
	}
	return boards
}

func SearchBoards(searchText string) (festructs.CBoardBatch, search.ScoreMap, error) {
	resp, err := search.Search(searchText)
	if err != nil {
		return []festructs.CompiledBoard{}, make(search.ScoreMap), err
	}
	b := festructs.CBoardBatch(findBoards(resp))
	b.SortByThreadsCount()
	return b, makeScoreMap(resp), nil
}

/*----------  Search & find users  ----------*/

func findUser(searchResult search.SearchResult) (festructs.CompiledUser, error) {
	searchId := searchResult.Id
	c := festructs.UserHeaderCarrier{}
	logging.Logf(3, "Single read happens in findUser>One")
	err := globals.KvInstance.One("Fingerprint", searchId.Fingerprint, &c)
	if err != nil {
		logging.Logf(1, "We could not find this user. SearchId: %v", searchId)
		return festructs.CompiledUser{}, err
	}
	return c.Users[0], nil
}

func findUsers(searchResults search.SearchResults) []festructs.CompiledUser {
	var users []festructs.CompiledUser

	for k := range searchResults.Results {
		if searchResults.Results[k].Id.EntityType != "User" {
			continue
		}
		p, err := findUser(searchResults.Results[k])
		if err != nil {
			continue
		}
		users = append(users, p)
	}
	return users
}

func SearchUsers(searchText string) (festructs.CUserBatch, search.ScoreMap, error) {
	resp, err := search.Search(searchText)
	if err != nil {
		return []festructs.CompiledUser{}, make(search.ScoreMap), err
	}
	return festructs.CUserBatch(findUsers(resp)), makeScoreMap(resp), nil
}

/*----------  Search & find content (posts + threads)  ----------*/

/*
Content is a composite type. This is more efficient than doing two searches, one for threads and one for posts, and then merging them together.
*/
func findContent(searchResults search.SearchResults) ([]festructs.CompiledPost, []festructs.CompiledThread) {
	var posts []festructs.CompiledPost

	var threads []festructs.CompiledThread

	for k := range searchResults.Results {
		if searchResults.Results[k].Id.EntityType != "Post" && searchResults.Results[k].Id.EntityType != "Thread" {
			continue
		}
		if searchResults.Results[k].Id.EntityType == "Post" {
			p, err := findPost(&searchResults.Results[k].Id)
			if err != nil {
				continue
			}
			posts = append(posts, p)
		}
		if searchResults.Results[k].Id.EntityType == "Thread" {
			p, err := findThread(&searchResults.Results[k].Id)
			if err != nil {
				continue
			}
			threads = append(threads, p)
		}
	}
	return posts, threads
}

func SearchContent(searchText string) (festructs.CPostBatch, festructs.CThreadBatch, search.ScoreMap, error) {
	resp, err := search.Search(searchText)
	if err != nil {
		return festructs.CPostBatch{}, festructs.CThreadBatch{}, make(search.ScoreMap), err
	}
	posts, threads := findContent(resp)
	return festructs.CPostBatch(posts), festructs.CThreadBatch(threads), makeScoreMap(resp), nil
}

/*----------  Make score map  ----------*/

func makeScoreMap(sr search.SearchResults) search.ScoreMap {
	sm := make(search.ScoreMap)
	for k := range sr.Results {
		sm[sr.Results[k].Id.Fingerprint] = sr.Results[k].Score
	}
	return sm
}

/*=======================================
=            New feed fetch            =
=======================================*/

func GetNewFeedContent(newFeedItems []beapiconsumer.NewFeedItem) ([]festructs.CompiledPost, []festructs.CompiledThread) {
	var posts []festructs.CompiledPost

	var threads []festructs.CompiledThread

	for k := range newFeedItems {
		if newFeedItems[k].EntityType != "Post" &&
			newFeedItems[k].EntityType != "Thread" &&
			newFeedItems[k].EntityType != "Vote" {
			continue
		}
		if newFeedItems[k].EntityType == "Post" {
			p, err := findPost(&newFeedItems[k])
			if err != nil {
				continue
			}
			posts = append(posts, p)
		}
		if newFeedItems[k].EntityType == "Thread" {
			p, err := findThread(&newFeedItems[k])
			if err != nil {
				continue
			}
			threads = append(threads, p)
		}
		if newFeedItems[k].EntityType == "Vote" {
			/*
				Votes are indicators for updates, we do not actually list votes in the new feed. However, if a user votes on the new feed, we need to be able to update the upvoted item after the vote is committed.

				Why? Imagine this, the user votes on the new feed. This vote is fully committed and in there.

				However, our new feed only listens to posts and threads. So a vote arriving in does not make new feed trigger a refresh on the target of that vote.

				That means even if the user votes on something, the new feed will not reflect that vote, since it still keeps the old copy of the entity before the vote, and it doesn't recognise that copy is now stale.

				So we need to determine if the target of this vote is a post or a thread, and update the appropriate item. Adding this here as a content update will remove the prior copy, which will leave the correctly updated new item as the only one.

				So we use the votes in this new feed as indicators on which additional posts and threads we need to update, not as directly visible votes shown.
			*/
			if newFeedItems[k].TargetFp == newFeedItems[k].ThreadFp {
				// This is a thread
				// We generate a surrogate to be able to search for this thread.
				threadSurrogateFeedItem := beapiconsumer.NewFeedItem{
					Fingerprint: newFeedItems[k].ThreadFp,
				}
				p, err := findThread(&threadSurrogateFeedItem)
				if err != nil {
					continue
				}
				threads = append(threads, p)
			} else {
				// This is a post
				postSurrogateFeedItem := beapiconsumer.NewFeedItem{
					Fingerprint: newFeedItems[k].TargetFp,
					ThreadFp:    newFeedItems[k].ThreadFp,
				}
				p, err := findPost(&postSurrogateFeedItem)
				if err != nil {
					continue
				}
				posts = append(posts, p)
			}
		}
	}
	return posts, threads
}

/*=====  End of New feed fetch  ======*/
