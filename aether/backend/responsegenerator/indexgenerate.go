// Backend > ResponseGenerator > IndexGenerate
// This file contains the index generation logic from given main entities.

package responsegenerator

import (
	// "fmt"
	"aether-core/aether/io/api"
	// "aether-core/aether/io/persistence"
	"aether-core/aether/services/globals"
	"aether-core/aether/services/logging"
	"aether-core/aether/services/toolbox"
	// "aether-core/aether/services/randomhashgen"
	// "encoding/json"
	// "errors"
	"fmt"
	// "github.com/davecgh/go-spew/spew"
	// "io/ioutil"
	// "os"
	// "strconv"
	// "strings"
	// "time"
)

func createBoardIndex(entity *api.Board, pageNum int) api.BoardIndex {
	var entityIndex api.BoardIndex
	entityIndex.Creation = entity.Creation
	entityIndex.Fingerprint = entity.GetFingerprint()
	entityIndex.LastUpdate = entity.LastUpdate
	entityIndex.PageNumber = pageNum
	entityIndex.EntityVersion = entity.EntityVersion
	return entityIndex
}

func createThreadIndex(entity *api.Thread, pageNum int) api.ThreadIndex {
	var entityIndex api.ThreadIndex
	entityIndex.Board = entity.Board
	entityIndex.Creation = entity.Creation
	entityIndex.Fingerprint = entity.GetFingerprint()
	entityIndex.PageNumber = pageNum
	entityIndex.EntityVersion = entity.EntityVersion
	return entityIndex
}

func createPostIndex(entity *api.Post, pageNum int) api.PostIndex {
	var entityIndex api.PostIndex
	entityIndex.Board = entity.Board
	entityIndex.Thread = entity.Thread
	entityIndex.Creation = entity.Creation
	entityIndex.Fingerprint = entity.GetFingerprint()
	entityIndex.PageNumber = pageNum
	entityIndex.EntityVersion = entity.EntityVersion
	return entityIndex
}

func createVoteIndex(entity *api.Vote, pageNum int) api.VoteIndex {
	var entityIndex api.VoteIndex
	entityIndex.Board = entity.Board
	entityIndex.Thread = entity.Thread
	entityIndex.Target = entity.Target
	entityIndex.Creation = entity.Creation
	entityIndex.Fingerprint = entity.GetFingerprint()
	entityIndex.LastUpdate = entity.LastUpdate
	entityIndex.PageNumber = pageNum
	entityIndex.EntityVersion = entity.EntityVersion
	return entityIndex
}

func createKeyIndex(entity *api.Key, pageNum int) api.KeyIndex {
	var entityIndex api.KeyIndex
	entityIndex.Creation = entity.Creation
	entityIndex.Fingerprint = entity.GetFingerprint()
	entityIndex.LastUpdate = entity.LastUpdate
	entityIndex.PageNumber = pageNum
	entityIndex.EntityVersion = entity.EntityVersion
	return entityIndex
}

func createTruststateIndex(entity *api.Truststate, pageNum int) api.TruststateIndex {
	var entityIndex api.TruststateIndex
	entityIndex.Target = entity.Target
	entityIndex.Creation = entity.Creation
	entityIndex.Fingerprint = entity.GetFingerprint()
	entityIndex.LastUpdate = entity.LastUpdate
	entityIndex.PageNumber = pageNum
	entityIndex.EntityVersion = entity.EntityVersion
	return entityIndex
}

// createUnbakedIndexes creates the index variant of every entity in an api.Response, and puts it back inside one single container for all indexes.
func createUnbakedIndexes(fullData *[]api.Response) *api.Response {
	fd := *fullData
	var resp api.Response
	if len(fd) > 0 {
		for i, _ := range fd {
			// For each Api.Response page
			if len(fd[i].Boards) > 0 {
				for j, _ := range fd[i].Boards {
					entityIndex := createBoardIndex(&fd[i].Boards[j], i)
					resp.BoardIndexes = append(resp.BoardIndexes, entityIndex)
				}
			}
			if len(fd[i].Threads) > 0 {
				for j, _ := range fd[i].Threads {
					entityIndex := createThreadIndex(&fd[i].Threads[j], i)
					resp.ThreadIndexes = append(resp.ThreadIndexes, entityIndex)
				}
			}
			if len(fd[i].Posts) > 0 {
				for j, _ := range fd[i].Posts {
					entityIndex := createPostIndex(&fd[i].Posts[j], i)
					resp.PostIndexes = append(resp.PostIndexes, entityIndex)
				}
			}
			if len(fd[i].Votes) > 0 {
				for j, _ := range fd[i].Votes {
					entityIndex := createVoteIndex(&fd[i].Votes[j], i)
					resp.VoteIndexes = append(resp.VoteIndexes, entityIndex)
				}
			}
			// Addresses: Address doesn't have an index form. It is its own index.
			// Addresses are skipped here.
			if len(fd[i].Keys) > 0 {
				for j, _ := range fd[i].Keys {
					entityIndex := createKeyIndex(&fd[i].Keys[j], i)
					resp.KeyIndexes = append(resp.KeyIndexes, entityIndex)
				}
			}
			if len(fd[i].Truststates) > 0 {
				for j, _ := range fd[i].Truststates {
					entityIndex := createTruststateIndex(&fd[i].Truststates[j], i)
					resp.TruststateIndexes = append(resp.TruststateIndexes, entityIndex)
				}
			}
		}
	}
	return &resp
}

func bakeIndexes(indexPages *[]api.ApiResponse, entityCounts *[]api.EntityCount, filters *[]api.Filter, foldername string, isPOST bool, respType string, entityType string) {
	// Create directory
	if respType == "addresses" {
		return // addresses do not generate indexes.
	}
	var indexdir string
	protv := globals.BackendConfig.GetProtURLVersion()
	if isPOST {
		indexdir = fmt.Sprint(globals.BackendConfig.GetCachesDirectory(), "/", protv, "/responses/", foldername, "/index")
	} else {
		if respType == "addresses" {
			indexdir = fmt.Sprint(globals.BackendConfig.GetCachesDirectory(), "/", protv, "/", respType, "/", foldername, "/index")
		} else {
			indexdir = fmt.Sprint(globals.BackendConfig.GetCachesDirectory(), "/", protv, "/c0/", respType, "/", foldername, "/index")
		}
	}
	toolbox.CreatePath(indexdir)
	for key, val := range *indexPages {
		// Add counts
		val.Caching.EntityCounts = *entityCounts
		// Add filters
		if filters != nil {
			val.Filters = *filters
		}
		val.Entity = entityType
		val.Endpoint = fmt.Sprint(entityType, "_index")
		if isPOST {
			val.Endpoint = fmt.Sprint(val.Endpoint, "_post")
		}
		// Sign
		signingErr := val.CreateSignature(globals.BackendConfig.GetBackendKeyPair())
		if signingErr != nil {
			logging.Log(1, fmt.Sprintf("This page of the index of a multiple-page post response failed to be page-signed. Error: %#v Page: %#v\n", signingErr, val))
		}
		// Convert to JSON
		indexJsonResp, err := val.ToJSON()
		if err != nil {
			logging.Log(1, fmt.Sprintf("This page of the index of a multiple-page post response failed to convert to JSON. Error: %#v\n", err))
		}
		// Save to disk
		filename := fmt.Sprint(key, ".json")
		saveFileToDisk(indexJsonResp, indexdir, filename)
	}
}
