// Backend > ResponseGenerator > Splitter
// This file provides a set of functions that splits a chunk of data into pages that can be transmitted over the network.

package responsegenerator

import (
	// "fmt"
	"aether-core/aether/io/api"
	// "aether-core/aether/io/persistence"
	"aether-core/aether/services/globals"
	// "aether-core/aether/services/logging"
	// "aether-core/aether/services/randomhashgen"
	// "encoding/json"
	// "errors"
	// "fmt"
	// "github.com/davecgh/go-spew/spew"
	// "io/ioutil"
	// "os"
	// "strconv"
	// "strings"
	// "time"
)

/*
splitManifestToPages is split off from the main entity splitter because manifest is a two-level deep entity, unlike indexes and regular entities, which are one-level deep.

In other words, manifest structure is:

responsebody >

	posts_manifest >
	  page:0 >
	    m1, m2, m3
	  page:1 >
	    m4,m5,m6 ..

while others are:

responsebody >

	posts >
	  e1, e2, e3 ...

This means if we count the first level entity count as a page splitting gate, it's not gonna work. We need to count manifests themselves. We can do this based on the manifests and try to figure out which page to put each page:0 item, but that's going to be wonky. In the default case it might not matter, but in the case where somebody breaks the config in a way that the page entity counts are vastly higher than the manifest counts (i.e. entity page takes 60k items while manifest takes 30k items) it might cause weird manifest sizings.

So: let's do this right. This means, we'll be inserting an unbaked manifest carrier here, splitting the manifests first, and then creating the intermediate layer of pages, and the inserting those to the response. So if you have 60k entity page and 30k manifest pages, the entity page will just be split into two manifest pages, both named page:0.
*/
func splitManifestToPages(fullData *unbakedManifestCarrier) *[]api.Response {
	var pages []api.Response
	var entityTypes []string
	// Manifests
	if len(fullData.BoardManifests) > 0 {
		entityTypes = append(entityTypes, "boardmanifests")
	}
	if len(fullData.ThreadManifests) > 0 {
		entityTypes = append(entityTypes, "threadmanifests")
	}
	if len(fullData.PostManifests) > 0 {
		entityTypes = append(entityTypes, "postmanifests")
	}
	if len(fullData.VoteManifests) > 0 {
		entityTypes = append(entityTypes, "votemanifests")
	}
	if len(fullData.KeyManifests) > 0 {
		entityTypes = append(entityTypes, "keymanifests")
	}
	if len(fullData.TruststateManifests) > 0 {
		entityTypes = append(entityTypes, "truststatemanifests")
	}
	if len(fullData.AddressManifests) > 0 {
		entityTypes = append(entityTypes, "addressmanifests")
	}
	if len(fullData.BoardManifests) == 0 &&
		len(fullData.ThreadManifests) == 0 &&
		len(fullData.PostManifests) == 0 &&
		len(fullData.VoteManifests) == 0 &&
		len(fullData.KeyManifests) == 0 &&
		len(fullData.TruststateManifests) == 0 &&
		len(fullData.AddressManifests) == 0 {
		entityTypes = append(entityTypes, "blankpage")
		// Why? because we still want to generate a blank manifest page if there is nothing inside, to communicate that this cache is empty.
	}
	for i := range entityTypes {
		if entityTypes[i] == "boardmanifests" {
			dataSet := fullData.BoardManifests
			pageSize := globals.BackendConfig.GetEntityPageSizes().BoardManifests
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.BoardManifests = *constructManifestStructure(&pageData)
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "threadmanifests" {
			dataSet := fullData.ThreadManifests
			pageSize := globals.BackendConfig.GetEntityPageSizes().ThreadManifests
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.ThreadManifests = *constructManifestStructure(&pageData)
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "postmanifests" {
			dataSet := fullData.PostManifests
			pageSize := globals.BackendConfig.GetEntityPageSizes().PostManifests
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.PostManifests = *constructManifestStructure(&pageData)
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "votemanifests" {
			dataSet := fullData.VoteManifests
			pageSize := globals.BackendConfig.GetEntityPageSizes().VoteManifests
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.VoteManifests = *constructManifestStructure(&pageData)
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "keymanifests" {
			dataSet := fullData.KeyManifests
			pageSize := globals.BackendConfig.GetEntityPageSizes().KeyManifests
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.KeyManifests = *constructManifestStructure(&pageData)
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "truststatemanifests" {
			dataSet := fullData.TruststateManifests
			pageSize := globals.BackendConfig.GetEntityPageSizes().TruststateManifests
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.TruststateManifests = *constructManifestStructure(&pageData)
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "addressmanifests" {
			dataSet := fullData.AddressManifests
			pageSize := globals.BackendConfig.GetEntityPageSizes().AddressManifests
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.AddressManifests = *constructManifestStructure(&pageData)
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "blankpage" {
			pages = append(pages, api.Response{})
		}
	}
	return &pages
}

func splitEntitiesToPages(fullData *api.Response) *[]api.Response {
	var entityTypes []string
	// We do this check set below so that we don't run pagination logic on entity types that does not exist in this response. This is a bit awkward because there's no good way to iterate over fields of a struct.

	// Mainline entities
	if len(fullData.Boards) > 0 {
		entityTypes = append(entityTypes, "boards")
	}
	if len(fullData.Threads) > 0 {
		entityTypes = append(entityTypes, "threads")
	}
	if len(fullData.Posts) > 0 {
		entityTypes = append(entityTypes, "posts")
	}
	if len(fullData.Votes) > 0 {
		entityTypes = append(entityTypes, "votes")
	}
	if len(fullData.Addresses) > 0 {
		entityTypes = append(entityTypes, "addresses")
	}
	if len(fullData.Keys) > 0 {
		entityTypes = append(entityTypes, "keys")
	}
	if len(fullData.Truststates) > 0 {
		entityTypes = append(entityTypes, "truststates")
	}
	// Indexes
	if len(fullData.BoardIndexes) > 0 {
		entityTypes = append(entityTypes, "boardindexes")
	}
	if len(fullData.ThreadIndexes) > 0 {
		entityTypes = append(entityTypes, "threadindexes")
	}
	if len(fullData.PostIndexes) > 0 {
		entityTypes = append(entityTypes, "postindexes")
	}
	if len(fullData.VoteIndexes) > 0 {
		entityTypes = append(entityTypes, "voteindexes")
	}
	if len(fullData.KeyIndexes) > 0 {
		entityTypes = append(entityTypes, "keyindexes")
	}
	if len(fullData.TruststateIndexes) > 0 {
		entityTypes = append(entityTypes, "truststateindexes")
	}
	if len(fullData.AddressIndexes) > 0 {
		entityTypes = append(entityTypes, "addressindexes")
	}

	var pages []api.Response
	// This is a lot of copy paste. This is because there is no automatic conversion from []api.Boards being recognised as []api.Provable. Without that, I have to convert them explicitly to be able to put them into a map[string:struct] which is a lot of extra work - more work than copy paste.
	for i := range entityTypes {
		if entityTypes[i] == "boards" {
			dataSet := fullData.Boards
			pageSize := globals.BackendConfig.GetEntityPageSizes().Boards
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.Boards = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "threads" {
			dataSet := fullData.Threads
			pageSize := globals.BackendConfig.GetEntityPageSizes().Threads
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.Threads = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "posts" {
			dataSet := fullData.Posts
			pageSize := globals.BackendConfig.GetEntityPageSizes().Posts
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.Posts = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "votes" {
			dataSet := fullData.Votes
			pageSize := globals.BackendConfig.GetEntityPageSizes().Votes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.Votes = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "addresses" {
			dataSet := fullData.Addresses
			pageSize := globals.BackendConfig.GetEntityPageSizes().Addresses
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.Addresses = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "keys" {
			dataSet := fullData.Keys
			pageSize := globals.BackendConfig.GetEntityPageSizes().Keys
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.Keys = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "truststates" {
			dataSet := fullData.Truststates
			pageSize := globals.BackendConfig.GetEntityPageSizes().Truststates
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.Truststates = pageData
				pages = append(pages, page)
			}
		}
		// Index entities
		if entityTypes[i] == "boardindexes" {
			dataSet := fullData.BoardIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().BoardIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.BoardIndexes = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "threadindexes" {
			dataSet := fullData.ThreadIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().ThreadIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.ThreadIndexes = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "postindexes" {
			dataSet := fullData.PostIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().PostIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.PostIndexes = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "voteindexes" {
			dataSet := fullData.VoteIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().VoteIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.VoteIndexes = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "keyindexes" {
			dataSet := fullData.KeyIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().KeyIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.KeyIndexes = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "addressindexes" {
			dataSet := fullData.AddressIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().AddressIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.AddressIndexes = pageData
				pages = append(pages, page)
			}
		}
		if entityTypes[i] == "truststateindexes" {
			dataSet := fullData.TruststateIndexes
			pageSize := globals.BackendConfig.GetEntityPageSizes().TruststateIndexes
			numPages := len(dataSet)/pageSize + 1
			// The division above is floored.
			for i := 0; i < numPages; i++ {
				beg := i * pageSize
				var end int
				// This is to protect from 'slice bounds out of range'
				if (i+1)*pageSize > len(dataSet) {
					end = len(dataSet)
				} else {
					end = (i + 1) * pageSize
				}
				pageData := dataSet[beg:end]
				var page api.Response
				page.TruststateIndexes = pageData
				pages = append(pages, page)
			}
		}
	}
	if len(entityTypes) == 0 {
		// The result is empty
		var page api.Response
		pages = append(pages, page)
	}
	return &pages
}
