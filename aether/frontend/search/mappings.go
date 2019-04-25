// Frontend > Search > Mappings
// This file handles the mappings between the search index and the frontend payload.

package search

import (
	"aether-core/aether/services/logging"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/analysis/analyzer/simple"
	"github.com/blevesearch/bleve/analysis/analyzer/standard"
	bleveMapping "github.com/blevesearch/bleve/mapping"
)

// buildMappings builds the mappings of our content to the index database. This is where we decide which field we want to search, in which form, and which fields to ignore.
func buildMappings() bleveMapping.IndexMapping {
	coreMapping := bleve.NewIndexMapping()
	coreMapping.AddDocumentMapping("board", generateBoardMapping())
	coreMapping.AddDocumentMapping("thread", generateThreadMapping())
	coreMapping.AddDocumentMapping("post", generatePostMapping())
	coreMapping.AddDocumentMapping("user", generateUserMapping())
	coreMapping.DefaultMapping = bleve.NewDocumentDisabledMapping()
	// ^Disabled by default - if we see something we don't recognise, we don't index.
	return coreMapping
}

func makeFieldMapping(maptype, analyser string) *bleveMapping.FieldMapping {
	switch maptype {
	case "text":
		textFieldMapping := bleve.NewTextFieldMapping()
		textFieldMapping.Store = false
		switch analyser {
		case "simple":
			textFieldMapping.Analyzer = simple.Name
		case "standard":
			textFieldMapping.Analyzer = standard.Name
		case "keyword":
			textFieldMapping.Analyzer = keyword.Name
		default:
			logging.LogCrashf("This requested analyser wasn't understood. %v", analyser)
		}
		return textFieldMapping
	case "numeric":
		numericFieldMapping := bleve.NewNumericFieldMapping()
		numericFieldMapping.Store = false
		return numericFieldMapping
	default:
		logging.LogCrashf("This requested mapping type wasn't understood. %v", maptype)
		return nil
	}
}

/*
	Board: Index:
	- Fingerprint (simple analyser)
	- Name (simple)
	- Description (standard analyser)
	- Creation
	- LastUpdate
	- ThreadsCount
	- UserCount
*/
func generateBoardMapping() *bleveMapping.DocumentMapping {
	mapping := bleve.NewDocumentStaticMapping()
	mapping.AddFieldMappingsAt("Fingerprint", makeFieldMapping("text", "keyword"))
	mapping.AddFieldMappingsAt("Name", makeFieldMapping("text", "simple"))
	mapping.AddFieldMappingsAt("Description", makeFieldMapping("text", "standard"))
	// mapping.AddFieldMappingsAt("Creation", makeFieldMapping("numeric", ""))
	// mapping.AddFieldMappingsAt("LastUpdate", makeFieldMapping("numeric", ""))
	// mapping.AddFieldMappingsAt("ThreadsCount", makeFieldMapping("numeric", ""))
	// mapping.AddFieldMappingsAt("UserCount", makeFieldMapping("numeric", ""))
	return mapping
}

/*
	Thread: Index:
	- Fingerprint (simple analyser)
	- Board (simple)
	- Name (simple)
	- Body (standard analyser)
	- Link (simple)
	- Creation
	- LastUpdate
	- PostsCount
	- Score
*/
func generateThreadMapping() *bleveMapping.DocumentMapping {
	mapping := bleve.NewDocumentStaticMapping()

	mapping.AddFieldMappingsAt("Fingerprint", makeFieldMapping("text", "keyword"))
	// mapping.AddFieldMappingsAt("Board", makeFieldMapping("text", "keyword"))
	mapping.AddFieldMappingsAt("Name", makeFieldMapping("text", "simple"))
	mapping.AddFieldMappingsAt("Body", makeFieldMapping("text", "standard"))
	mapping.AddFieldMappingsAt("Link", makeFieldMapping("text", "keyword"))
	// mapping.AddFieldMappingsAt("Creation", makeFieldMapping("numeric", ""))
	// mapping.AddFieldMappingsAt("LastUpdate", makeFieldMapping("numeric", ""))
	// mapping.AddFieldMappingsAt("PostsCount", makeFieldMapping("numeric", ""))
	// mapping.AddFieldMappingsAt("Score", makeFieldMapping("numeric", ""))
	return mapping
}

/*
	Post: Index:
	- Fingerprint (simple analyser)
	- Board (simple)
	- Thread (simple)
	- Parent (simple)
	- Body (standard analyser)
	- Creation
	- LastUpdate
*/
func generatePostMapping() *bleveMapping.DocumentMapping {
	mapping := bleve.NewDocumentStaticMapping()

	mapping.AddFieldMappingsAt("Fingerprint", makeFieldMapping("text", "keyword"))
	// mapping.AddFieldMappingsAt("Board", makeFieldMapping("text", "keyword"))
	// mapping.AddFieldMappingsAt("Thread", makeFieldMapping("text", "keyword"))
	// mapping.AddFieldMappingsAt("Parent", makeFieldMapping("text", "keyword"))
	mapping.AddFieldMappingsAt("Body", makeFieldMapping("text", "standard"))
	// mapping.AddFieldMappingsAt("Creation", makeFieldMapping("numeric", ""))
	// mapping.AddFieldMappingsAt("LastUpdate", makeFieldMapping("numeric", ""))
	return mapping
}

/*
	User: Index:
	- Fingerprint (simple analyser)
	- NonCanonicalName (simple)
	- Info (standard analyser)
	- Creation
	- LastUpdate
	- LastRefreshed
*/
func generateUserMapping() *bleveMapping.DocumentMapping {
	mapping := bleve.NewDocumentStaticMapping()
	mapping.AddFieldMappingsAt("Fingerprint", makeFieldMapping("text", "keyword"))
	mapping.AddFieldMappingsAt("NonCanonicalName", makeFieldMapping("text", "simple"))
	mapping.AddFieldMappingsAt("Info", makeFieldMapping("text", "standard"))
	// mapping.AddFieldMappingsAt("Creation", makeFieldMapping("numeric", ""))
	// mapping.AddFieldMappingsAt("LastUpdate", makeFieldMapping("numeric", ""))
	// mapping.AddFieldMappingsAt("LastRefreshed", makeFieldMapping("numeric", ""))

	cusmapping := bleve.NewDocumentStaticMapping()
	cusmapping.AddFieldMappingsAt("CanonicalName", makeFieldMapping("text", "simple"))

	mapping.AddSubDocumentMapping("CompiledUserSignals", cusmapping)
	return mapping
}
