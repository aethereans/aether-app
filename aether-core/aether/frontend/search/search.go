// Frontend > Search
// This package handles the indexing and full text search in the frontend key/value store.

package search

import (
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"github.com/blevesearch/bleve"
	// bleveMapping "github.com/blevesearch/bleve/mapping"
	// "github.com/davecgh/go-spew/spew"
	"github.com/json-iterator/go"
	"os"
	"path/filepath"
)

var (
	json  = jsoniter.ConfigCompatibleWithStandardLibrary // declare, assign
	index bleve.Index                                    // declare only. (i.e. missing = is not a typo)
)

/*----------  Utility functions and data containers  ----------*/

// Search Id structure

type SearchId struct {
	EntityType  string
	BoardFp     string
	ThreadFp    string
	ParentFp    string
	Fingerprint string
	UserFp      string
}

// Make it satisfy Searchable interface
func (o *SearchId) GetEntityType() string {
	return o.EntityType
}
func (o *SearchId) GetBoardFp() string {
	return o.BoardFp
}
func (o *SearchId) GetThreadFp() string {
	return o.ThreadFp
}
func (o *SearchId) GetParentFp() string {
	return o.ParentFp
}
func (o *SearchId) GetFingerprint() string {
	return o.Fingerprint
}
func (o *SearchId) GetUserFp() string {
	return o.UserFp
}

type SearchResult struct {
	Id    SearchId
	Score float64
}

type SearchResults struct {
	Query   string
	Results []SearchResult
}

type ScoreMap map[string]float64

func MakeSearchId(entityType, boardfp, threadfp, parentfp, fingerprint, userfp string) (string, error) {
	idStruct := SearchId{
		EntityType:  entityType, // Board, Thread, Post, User
		BoardFp:     boardfp,
		ThreadFp:    threadfp,
		ParentFp:    parentfp,
		Fingerprint: fingerprint,
		UserFp:      userfp,
	}
	idByte, err := json.Marshal(idStruct)
	if err != nil {
		return "", err
	}
	return string(idByte), nil
}

func IndexExists() bool {
	idir := filepath.Join(globals.FrontendConfig.GetUserDirectory(), "frontend")
	iloc := filepath.Join(idir, "searchindex")
	if _, err := os.Stat(iloc); !os.IsNotExist(err) {
		return true
	}
	return false
}

func OpenIndex() {
	idir := filepath.Join(globals.FrontendConfig.GetUserDirectory(), "frontend")
	iloc := filepath.Join(idir, "searchindex")
	idx, err := bleve.Open(iloc)
	if err != nil {
		if err == bleve.ErrorIndexPathDoesNotExist {
			idx = createNewIndex(iloc)
		} else {
			logging.LogCrashf("Bleve Index Open errored out. Err: %v", err)
		}
	}
	index = idx
}

func createNewIndex(iloc string) bleve.Index {
	mapping := buildMappings()
	index, err := bleve.New(iloc, mapping)
	if err != nil {
		logging.LogCrash(err)
	}
	return index
}

func CloseIndex() {
	index.Close()
}

func Index(id string, data interface{}) error {
	return index.Index(id, data)
}

func Delete(id string) error {
	return index.Delete(id)
}

func NewBatch() *bleve.Batch {
	return index.NewBatch()
}

func CommitBatch(b *bleve.Batch) error {
	return index.Batch(b)
}

/*----------  Search  ----------*/

/*
Search returns IDs (locators) for the actual results. It does not return the results themselves. You still need to find those records from the KVStore yourself using the ID. The ID is a struct that carries enough information to find pretty much anything fast, though.

HEADS UP: Search normally does have a type called entityType string that it receives, but so far, our search does not actually care about the entity type we want. It just returns the whole results, and then we filter it at a higher level of the stack.

Whenever we become sophisticated enough to do the scanning on the entity type, it should be added back.
*/
func Search(searchText string) (SearchResults, error) {
	query := bleve.NewMatchQuery(searchText)
	search := bleve.NewSearchRequest(query)
	search.Size = 10000
	results, err := index.Search(search)
	if err != nil {
		logging.Logf(1, "Search encountered an error. Err: %v", err)
		return SearchResults{}, err
	}
	srs := SearchResults{Query: searchText}
	for k, _ := range results.Hits {
		resultId := SearchId{}
		err := json.Unmarshal([]byte(results.Hits[k].ID), &resultId)
		if err != nil {
			logging.Logf(1, "This search result could not be unmarshaled into a result id struct. Result: %v", results.Hits[k].ID)
			return srs, err
		}
		sr := SearchResult{
			Id:    resultId,
			Score: results.Hits[k].Score,
		}
		srs.Results = append(srs.Results, sr)
	}
	logging.Logf(1, "Search results: %#v", srs)
	return srs, nil
}
