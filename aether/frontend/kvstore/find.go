// Frontend > KeyValueStore > Find
// This file takes some search IDs and it converts them to actual results we can deliver to the user.

package kvstore

import (
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

func findPost(searchResult search.SearchResult) (festructs.CompiledPost, error) {
	searchId := searchResult.Id
	c := festructs.ThreadCarrier{}
	err := globals.KvInstance.One("Fingerprint", searchId.ThreadFp, &c)
	if err != nil {
		logging.Logf(1, "We could not find the thread of this post. SearchId: %v", searchId)
		return festructs.CompiledPost{}, err
	}
	i := c.Posts.Find(searchId.Fingerprint)
	return c.Posts[i], nil
}

func findPosts(searchResults search.SearchResults) []festructs.CompiledPost {
	posts := []festructs.CompiledPost{}
	for k, _ := range searchResults.Results {
		if searchResults.Results[k].Id.EntityType != "Post" {
			continue
		}
		p, err := findPost(searchResults.Results[k])
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

func findThread(searchResult search.SearchResult) (festructs.CompiledThread, error) {
	searchId := searchResult.Id
	c := festructs.ThreadCarrier{}
	err := globals.KvInstance.One("Fingerprint", searchId.Fingerprint, &c)
	if err != nil {
		logging.Logf(1, "We could not find this thread. SearchId: %v", searchId)
		return festructs.CompiledThread{}, err
	}
	return c.Threads[0], nil
}

func findThreads(searchResults search.SearchResults) []festructs.CompiledThread {
	threads := []festructs.CompiledThread{}
	for k, _ := range searchResults.Results {
		if searchResults.Results[k].Id.EntityType != "Thread" {
			continue
		}
		p, err := findThread(searchResults.Results[k])
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
	err := globals.KvInstance.One("Fingerprint", searchId.Fingerprint, &c)
	if err != nil {
		logging.Logf(1, "We could not find this board. SearchId: %v", searchId)
		return festructs.CompiledBoard{}, err
	}
	return c.Boards[0], nil
}

func findBoards(searchResults search.SearchResults) []festructs.CompiledBoard {
	boards := []festructs.CompiledBoard{}
	for k, _ := range searchResults.Results {
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
	err := globals.KvInstance.One("Fingerprint", searchId.Fingerprint, &c)
	if err != nil {
		logging.Logf(1, "We could not find this user. SearchId: %v", searchId)
		return festructs.CompiledUser{}, err
	}
	return c.Users[0], nil
}

func findUsers(searchResults search.SearchResults) []festructs.CompiledUser {
	users := []festructs.CompiledUser{}
	for k, _ := range searchResults.Results {
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
	posts := []festructs.CompiledPost{}
	threads := []festructs.CompiledThread{}
	for k, _ := range searchResults.Results {
		if searchResults.Results[k].Id.EntityType != "Post" && searchResults.Results[k].Id.EntityType != "Thread" {
			continue
		}
		if searchResults.Results[k].Id.EntityType == "Post" {
			p, err := findPost(searchResults.Results[k])
			if err != nil {
				continue
			}
			posts = append(posts, p)
		}
		if searchResults.Results[k].Id.EntityType == "Thread" {
			p, err := findThread(searchResults.Results[k])
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
	for k, _ := range sr.Results {
		sm[sr.Results[k].Id.Fingerprint] = sr.Results[k].Score
	}
	return sm
}
